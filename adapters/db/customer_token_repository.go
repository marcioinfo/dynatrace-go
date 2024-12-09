package db

import (
	"context"
	"database/sql"
	customertoken "payment-layer-card-api/entities/customer_token"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CustomerTokenRepository struct {
	db *sql.DB
}

func NewCustomerTokenRepository(db *sql.DB) *CustomerTokenRepository {
	return &CustomerTokenRepository{
		db: db,
	}
}

func (r *CustomerTokenRepository) Insert(ctx context.Context, customerToken *customertoken.CustomerToken) *errors.ErrorOutput {
	err := r.db.QueryRowContext(ctx, `
         INSERT INTO customer_tokens (
                                        id,
                                        customer_id,
                                        customer_token,
                                        gateway
                                     )
                                     VALUES
                                     (
                                        $1,
                                        $2,
                                        $3,
                                        $4
                                     ) RETURNING id`,
		customerToken.ID,
		customerToken.CustomerId,
		customerToken.CustomerToken,
		customerToken.Gateway,
	).Scan(&customerToken.ID)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}

func (r *CustomerTokenRepository) GetByCustomerID(ctx context.Context, customerID string) ([]*customertoken.CustomerToken, *errors.ErrorOutput) {
	rows, err := r.db.QueryContext(ctx, `
         SELECT id,
                customer_id,
                customer_token,
                gateway,
                created_at,
                updated_at
           FROM customer_tokens
          WHERE customer_id = $1`,
		customerID,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	defer rows.Close()

	var customerTokens []*customertoken.CustomerToken

	for rows.Next() {
		customerToken := &customertoken.CustomerToken{}

		if err := rows.Scan(&customerToken.ID,
			&customerToken.CustomerId,
			&customerToken.CustomerToken,
			&customerToken.Gateway,
			&customerToken.CreatedAt,
			&customerToken.UpdatedAt); err != nil {
			return nil, errors.NewError(errors.InternalServerError, err)
		}

		customerTokens = append(customerTokens, customerToken)
	}

	return customerTokens, nil
}

func (r *CustomerTokenRepository) GetByCustomerIDAndGateway(ctx context.Context, customerId string, gateway string) (*customertoken.CustomerToken, *errors.ErrorOutput) {
	customerToken := &customertoken.CustomerToken{}

	err := r.db.QueryRowContext(ctx, `
         SELECT id,
                customer_id,
                customer_token,
                gateway,
                created_at,
                updated_at
           FROM customer_tokens
          WHERE customer_id = $1
            AND gateway     = $2`,
		customerId,
		gateway).Scan(&customerToken.ID,
		&customerToken.CustomerId,
		&customerToken.CustomerToken,
		&customerToken.Gateway,
		&customerToken.CreatedAt,
		&customerToken.UpdatedAt)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	return customerToken, nil
}
