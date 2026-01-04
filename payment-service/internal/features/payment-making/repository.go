package paymentmaking

import (
	"context"
	"payment-service/internal/db"
	"payment-service/internal/features/common"
	"payment-service/internal/models"

	"github.com/jackc/pgx/v5"
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
	GetTransactionByID(
		ctx context.Context,
		id string,
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
	RETURNING id, amount, currency, status, created_at;
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
		&result.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return r.GetTransactionByRef(ctx, tx.Ref)
	}

	if err != nil {
		return nil, common.ErrInternal
	}

	return &result, nil
}

func (r *repository) GetTransactionByRef(
	ctx context.Context,
	ref string,
) (*models.Transaction, error) {

	query := `
	SELECT id, amount, currency, status, created_at
	FROM transactions
	WHERE reference = $1
	`

	var tx models.Transaction

	err := r.db.ConnPool.QueryRow(ctx, query, ref).Scan(
		&tx.ID,
		&tx.Amount,
		&tx.Currency,
		&tx.Status,
		&tx.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, common.ErrNotFound
	}

	if err != nil {
		return nil, common.ErrInternal
	}
	return &tx, nil
}

func (r *repository) GetTransactionByID(
	ctx context.Context,
	id string,
) (*models.Transaction, error) {

	query := `
	SELECT id, amount, currency, reference, status, created_at
	FROM transactions
	WHERE id = $1
	`

	var tx models.Transaction

	err := r.db.ConnPool.QueryRow(ctx, query, id).Scan(
		&tx.ID,
		&tx.Amount,
		&tx.Currency,
		&tx.Ref,
		&tx.Status,
		&tx.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, common.ErrNotFound
	}

	if err != nil {
		return nil, common.ErrInternal
	}
	return &tx, nil
}
