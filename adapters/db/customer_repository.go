package db

import (
	"context"
	"database/sql"
	"payment-layer-card-api/entities/customers"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		DB: db,
	}
}

func (r *CustomerRepository) Insert(ctx context.Context, customer *customers.Customer) *errors.ErrorOutput {
	err := r.DB.QueryRowContext(ctx,
		`INSERT INTO customers (id, name, document, birth_date, email, phone, gender) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		customer.ID,
		customer.Name,
		customer.Document,
		customer.BirthDate,
		customer.Email,
		customer.Phone,
		customer.Gender,
	).Scan(&customer.ID)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}

func (r *CustomerRepository) GetByID(ctx context.Context, id string) (*customers.Customer, *errors.ErrorOutput) {
	customer := &customers.Customer{}
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, name, document, birth_date,email, phone, gender, created_at, updated_at  FROM customers WHERE id = $1 and deleted_at IS NULL`,
		id,
	).Scan(&customer.ID,
		&customer.Name,
		&customer.Document,
		&customer.BirthDate,
		&customer.Email,
		&customer.Phone,
		&customer.Gender,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	return customer, nil
}

func (r *CustomerRepository) GetCustomerByDocument(ctx context.Context, document string, serviceID string) (*customers.Customer, *errors.ErrorOutput) {
	customer := &customers.Customer{}

	var birthDate, createdAt, updatedAt sql.NullTime

	err := r.DB.QueryRowContext(ctx,
		`SELECT c.id, cs.name, cs.document, cs.birth_date, cs.email, cs.phone, cs.gender, cs.created_at, cs.updated_at
		FROM customers c
		JOIN customer_service cs ON c.id = cs.customer_id
		WHERE c.document = $1 AND cs.service_id = $2 AND c.deleted_at IS NULL`,
		document, serviceID,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Document,
		&birthDate,
		&customer.Email,
		&customer.Phone,
		&customer.Gender,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	if birthDate.Valid {
		customer.BirthDate = datetypes.CustomDate(birthDate.Time)
	}
	if createdAt.Valid {
		customer.CreatedAt = datetypes.CustomDateTime(createdAt.Time)
	}
	if updatedAt.Valid {
		customer.UpdatedAt = datetypes.CustomDateTime(updatedAt.Time)
	}

	return customer, nil
}

func (r *CustomerRepository) GetCustomerByEmail(ctx context.Context, email string) (*customers.Customer, *errors.ErrorOutput) {
	customer := &customers.Customer{}

	err := r.DB.QueryRowContext(ctx,
		`SELECT id, name, document, birth_date,email, phone, gender, created_at, updated_at  FROM customers WHERE email = $1 and deleted_at IS NULL`,
		email,
	).Scan(&customer.ID,
		&customer.Name,
		&customer.Document,
		&customer.BirthDate,
		&customer.Email,
		&customer.Phone,
		&customer.Gender,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	return customer, nil
}

func (r *CustomerRepository) Update(ctx context.Context, customer *customers.Customer) *errors.ErrorOutput {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE customers SET name = $1, document = $2, birth_date = $3, email = $4, phone= $5, gender= $6, updated_at = now() WHERE id = $7`,
		customer.Name,
		customer.Document,
		customer.BirthDate,
		customer.Email,
		customer.Phone,
		customer.Gender,
		customer.ID,
	)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}
