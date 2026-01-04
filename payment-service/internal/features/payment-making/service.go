package paymentmaking

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	broker "payment-service/internal/broker/rabbitmq"
	"payment-service/internal/features/common"
	"payment-service/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Service interface {
	CreatePayment(
		context.Context,
		CreatePaymentRequest,
	) (*CreatePaymentResponse, error)
	GetPaymentStatus(
		context.Context,
		string,
	) (*GetPaymentStatusResponse, error)
}

type service struct {
	rmq  *broker.RabbitMQ
	repo Repository
}

func NewService(r Repository, b *broker.RabbitMQ) Service {
	return &service{
		repo: r,
		rmq:  b,
	}
}

func (s *service) CreatePayment(
	ctx context.Context,
	req CreatePaymentRequest,
) (*CreatePaymentResponse, error) {
	parsedAmount, err := parseAmountToMinorUnits(req.Amount)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}
	transaction, err := s.repo.CreateOrGetTransaction(
		ctx,
		&models.Transaction{
			ID:        uuid.NewString(),
			Amount:    parsedAmount,
			Currency:  req.Currency,
			Ref:       req.Ref,
			Status:    models.TransactionPending,
			CreatedAt: time.Now().UTC(),
		},
	)
	if err != nil {
		return nil, err
	}

	if transaction.Status == models.TransactionPending {
		payload, _ := json.Marshal(broker.NewPayment{
			PaymentId: transaction.ID,
		})

		err = s.rmq.Publish(broker.Message{
			Exchange: broker.PaymentExchangeName,
			Topic:    broker.PaymentExchangeRoutingKey,
			Data:     payload,
		})
		if err != nil {
			return nil, echo.NewHTTPError(
				http.StatusInternalServerError,
				err.Error(),
			)
		}
		log.Println("successfuly published message")
	}

	return &CreatePaymentResponse{
		PaymentId: transaction.ID,
		Status:    transaction.Status,
	}, nil
}

func (s *service) GetPaymentStatus(ctx context.Context, id string) (*GetPaymentStatusResponse, error) {
	transaction, err := s.repo.GetTransactionByID(ctx, id)
	if err != nil {
		switch err {
		case common.ErrNotFound:
			return nil, echo.NewHTTPError(
				http.StatusNotFound,
				err.Error(),
			)
		default:
			return nil, echo.NewHTTPError(
				http.StatusInternalServerError,
				err.Error(),
			)
		}
	}

	return &GetPaymentStatusResponse{
		Amount:    transaction.DisplayAmount(),
		Currency:  transaction.Currency,
		Ref:       transaction.Ref,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
	}, nil
}
