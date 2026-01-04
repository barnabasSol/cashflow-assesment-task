package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"payment-worker/internal/config"
	"payment-worker/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) paymentHandler(
	p NewPayment,
	msg *amqp.Delivery,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := r.db.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		msg.Reject(false)
		log.Println("failed to begin transaction:", err)
		return
	}
	defer tx.Rollback(ctx)

	var transaction models.Transaction
	query := `
		SELECT id, amount, currency, reference, status, created_at
		FROM transactions
		WHERE id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, query, p.PaymentId).Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Ref,
		&transaction.Status,
		&transaction.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			msg.Nack(false, false)
		} else {
			log.Println(err)
			msg.Reject(false)
		}
		return
	}

	newStatus := RandomStatus()

	updateQuery := `
		UPDATE transactions
		SET status = $1
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateQuery, newStatus, p.PaymentId)
	if err != nil {
		msg.Reject(false)
		log.Println("failed to update transaction:", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		msg.Reject(false)
		log.Println("failed to commit transaction:", err)
		return
	}

	msg.Ack(false)
	log.Println("processed payment successfully:", transaction.ID, "new status:", newStatus)
}

func (r *RabbitMQ) ListenForNewPayment(msgs <-chan amqp.Delivery) {
	log.Println("payment listening...")

	for msg := range msgs {
		log.Printf("Received [%s]: %s", msg.RoutingKey, msg.Body)

		switch msg.RoutingKey {
		case PaymentExchangeRoutingKey:
			var payload NewPayment
			decoder := json.NewDecoder(bytes.NewReader(msg.Body))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&payload); err != nil {
				msg.Nack(false, false)
				log.Println("failed to unmarshal payload:", err)
				continue
			}
			retries := getRetryCount(&msg)
			if int32(retries) > config.GetEnvInt32("DLX_MAX_RETRY", 5) {
				log.Printf(
					"CRITICAL: Message %s failed 5 times. Dropping to prevent loop.",
					payload.PaymentId,
				)
				//Ack to stop the DLX loop
				msg.Ack(false)
				continue
			}
			r.paymentHandler(payload, &msg)

		default:
			log.Printf("unknown payment event: %s", msg.RoutingKey)
			msg.Nack(false, false)
		}
	}
}

func RandomStatus() int {
	if rand.Intn(2) == 0 {
		return -1
	}
	return 1
}
