package paymentmaking

import (
	"context"
	"net/http"
	broker "payment-service/internal/broker/rabbitmq"
	"payment-service/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Service interface {
	CreatePayment(context.Context, CreatePaymentRequest) (*CreatePaymentResponse, error)
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

	return &CreatePaymentResponse{
		PaymentId: transaction.ID,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
	}, nil
}
