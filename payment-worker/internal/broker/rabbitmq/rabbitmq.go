package broker

import (
	"payment-worker/internal/config"
	"payment-worker/internal/db"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	db   *db.Postgres
	conn *amqp.Connection
	ch   *amqp.Channel
}

func InitRabbitMQ(db *db.Postgres) (*RabbitMQ, error) {
	uri := config.GetSecret("RMQ_URL")
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = NewPaymentExchange(ch)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
		db:   db,
	}, nil
}

func (r *RabbitMQ) Publish(msg Message) error {
	return r.ch.Publish(
		msg.Exchange,
		msg.Topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Data,
		},
	)
}
func (r *RabbitMQ) Close() {
	if r.ch != nil {
		_ = r.ch.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
