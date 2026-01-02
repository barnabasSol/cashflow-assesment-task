package paymentmaking

import (
	"context"
	"net/http"
	"payment-service/internal/db"
	"payment-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	CreateOrGetTransaction(
		ctx context.Context,
		tx *models.Transaction,
	) (*models.Transaction, error)
	GetTransactionByRef(
		ctx context.Context,
		ref string,
	) (*models.Transaction, error)
}

type repository struct {
	db *db.Postgres
}

func NewRepository(db *db.Postgres) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateOrGetTransaction(
	ctx context.Context,
	tx *models.Transaction,
) (*models.Transaction, error) {

	query := `
	INSERT INTO transactions (id, reference, amount, currency, status)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (reference) DO NOTHING
	RETURNING id, amount, currency, status;
	`

	var result models.Transaction

	err := r.db.ConnPool.QueryRow(
		ctx,
		query,
		tx.ID,
		tx.Ref,
		tx.Amount,
		tx.Currency,
		tx.Status,
	).Scan(
		&result.ID,
		&result.Amount,
		&result.Currency,
		&result.Status,
	)

	if err == pgx.ErrNoRows {
		return r.GetTransactionByRef(ctx, tx.Ref)
	}

	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return &result, nil
}

func (r *repository) GetTransactionByRef(
	ctx context.Context,
	ref string,
) (*models.Transaction, error) {

	query := `
	SELECT id, amount, currency, status
	FROM transactions
	WHERE reference = $1
	`

	var tx models.Transaction

	err := r.db.ConnPool.QueryRow(ctx, query, ref).Scan(
		&tx.ID,
		&tx.Amount,
		&tx.Currency,
		&tx.Status,
	)

	if err == pgx.ErrNoRows {
		return nil, echo.NewHTTPError(
			http.StatusNotFound,
			"transaction doesn't exist",
		)
	}

	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return &tx, nil
}
