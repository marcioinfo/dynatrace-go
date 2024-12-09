package db

import (
	"context"
	"database/sql"
	"payment-layer-card-api/entities/cards"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CardRepository struct {
	DB *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{
		DB: db,
	}
}

func (r *CardRepository) Insert(ctx context.Context, card *cards.Card) *errors.ErrorOutput {
	err := r.DB.QueryRowContext(ctx,
		`INSERT INTO cards (id, customer_id, holder, brand, fingerprint, first_digits, last_digits, expire_month, expire_year) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		card.ID,
		card.CustomerID,
		card.Holder,
		card.Brand,
		card.Fingerprint,
		card.FirstDigits,
		card.LastDigits,
		card.ExpMonth,
		card.ExpYear,
	).Scan(&card.ID)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}

func (r *CardRepository) GetByID(ctx context.Context, id string) (*cards.Card, *errors.ErrorOutput) {
	card := &cards.Card{}

	err := r.DB.QueryRowContext(ctx,
		`SELECT id, customer_id, holder, brand, first_digits, last_digits, expire_month, expire_year, created_at, updated_at  FROM cards WHERE id = $1 and deleted_at IS NULL`,
		id,
	).Scan(&card.ID,
		&card.CustomerID,
		&card.Holder,
		&card.Brand,
		&card.FirstDigits,
		&card.LastDigits,
		&card.ExpMonth,
		&card.ExpYear,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}
	return card, nil
}

func (r *CardRepository) GetByIDAndCustomerID(ctx context.Context, id string, customerId string) (*cards.Card, *errors.ErrorOutput) {
	card := &cards.Card{}

	err := r.DB.QueryRowContext(ctx,
		`SELECT id, customer_id, holder, brand, first_digits, last_digits, expire_month, expire_year, created_at, updated_at FROM cards WHERE id = $1 and customer_id = $2 and deleted_at IS NULL`,
		id,
		customerId,
	).Scan(&card.ID,
		&card.CustomerID,
		&card.Holder,
		&card.Brand,
		&card.FirstDigits,
		&card.LastDigits,
		&card.ExpMonth,
		&card.ExpYear,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	return card, nil
}

func (r *CardRepository) GetCardsByCustomerID(ctx context.Context, customerID string) ([]*cards.Card, *errors.ErrorOutput) {
	cardsArray := []*cards.Card{}

	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, customer_id, holder, brand, first_digits, last_digits, expire_month, expire_year, created_at, updated_at  FROM cards WHERE customer_id = $1 and deleted_at IS NULL`,
		customerID,
	)
	if err != nil {
		return nil, errors.NewError(errors.InternalServerError, err)
	}
	defer rows.Close()

	for rows.Next() {
		card := &cards.Card{}

		err = rows.Scan(
			&card.ID,
			&card.CustomerID,
			&card.Holder,
			&card.Brand,
			&card.FirstDigits,
			&card.LastDigits,
			&card.ExpMonth,
			&card.ExpYear,
			&card.CreatedAt,
			&card.UpdatedAt,
		)

		if err != nil {
			return nil, errors.NewError(errors.InternalServerError, err)
		}

		cardsArray = append(cardsArray, card)
	}

	return cardsArray, nil
}

func (r *CardRepository) GetCardByFingerprint(ctx context.Context, fingerprint string) (*cards.Card, *errors.ErrorOutput) {
	card := &cards.Card{}

	err := r.DB.QueryRowContext(ctx, `
                SELECT id,
                       customer_id,
                       holder,
                       brand,
                       fingerprint,
                       first_digits,
                       last_digits,
                       expire_month,
                       expire_year,
                       created_at,
                       updated_at
                  FROM cards
                 WHERE fingerprint = $1
                   AND deleted_at IS NULL`,
		fingerprint,
	).Scan(&card.ID,
		&card.CustomerID,
		&card.Holder,
		&card.Brand,
		&card.Fingerprint,
		&card.FirstDigits,
		&card.LastDigits,
		&card.ExpMonth,
		&card.ExpYear,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	return card, nil
}

func (cr *CardRepository) DeleteByID(ctx context.Context, id string) (string, *errors.ErrorOutput) {
	query := "UPDATE cards SET deleted_at = now() WHERE id = $1"
	_, err := cr.DB.ExecContext(ctx, query, id)
	if err != nil {
		return "", errors.NewError(errors.InternalServerError, err)
	}

	return id, nil
}
