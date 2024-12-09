package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	customerservice "payment-layer-card-api/entities/customer_service"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CustomerServiceRepository struct {
	DB *sql.DB
}

func NewCustomerServiceRepository(db *sql.DB) *CustomerServiceRepository {
	return &CustomerServiceRepository{
		DB: db,
	}
}

func (r *CustomerServiceRepository) Insert(ctx context.Context, customerService *customerservice.CustomerService) *errors.ErrorOutput {
	err := r.DB.QueryRowContext(ctx,
		`INSERT INTO customer_service (id, service_id, customer_id, name, document, birth_date, email, phone, gender) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		customerService.ID,
		customerService.ServiceID,
		customerService.CustomerID,
		customerService.Name,
		customerService.Document,
		customerService.BirthDate,
		customerService.Email,
		customerService.Phone,
		customerService.Gender,
	).Scan(&customerService.ID)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}

func (r *CustomerServiceRepository) GetByCustomerID(ctx context.Context, customerID string) ([]*customerservice.CustomerService, *errors.ErrorOutput) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, service_id, customer_id, name, document, birth_date, email, phone, gender  
		 FROM customer_service WHERE customer_id = $1 AND deleted_at IS NULL`,
		customerID,
	)

	if err != nil {
		return nil, errors.NewError(errors.InternalServerError, err)
	}
	defer rows.Close()

	var customerServices []*customerservice.CustomerService
	for rows.Next() {
		customerService := &customerservice.CustomerService{}
		if err := rows.Scan(
			&customerService.ID,
			&customerService.ServiceID,
			&customerService.CustomerID,
			&customerService.Name,
			&customerService.Document,
			&customerService.BirthDate,
			&customerService.Email,
			&customerService.Phone,
			&customerService.Gender,
		); err != nil {
			return nil, errors.NewError(errors.InternalServerError, err)
		}
		customerServices = append(customerServices, customerService)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	if len(customerServices) == 0 {
		return nil, errors.NewError(errors.NotFoundError, nil, "customer service não encontrado")
	}

	return customerServices, nil
}

func (r *CustomerServiceRepository) GetByServiceID(ctx context.Context, serviceID string) ([]*customerservice.CustomerService, *errors.ErrorOutput) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, service_id, customer_id, name, document, birth_date, email, phone, gender, created_at, updated_at, deleted_at 
		 FROM customer_service WHERE service_id = $1 AND deleted_at IS NULL`,
		serviceID,
	)

	if err != nil {
		return nil, errors.NewError(errors.InternalServerError, err)
	}
	defer rows.Close()

	var customerServices []*customerservice.CustomerService
	for rows.Next() {
		customerService := &customerservice.CustomerService{}
		if err := rows.Scan(
			&customerService.ID,
			&customerService.ServiceID,
			&customerService.CustomerID,
			&customerService.Name,
			&customerService.Document,
			&customerService.BirthDate,
			&customerService.Email,
			&customerService.Phone,
			&customerService.Gender,
			&customerService.CreatedAt,
			&customerService.UpdatedAt,
		); err != nil {
			return nil, errors.NewError(errors.InternalServerError, err)
		}
		customerServices = append(customerServices, customerService)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	if len(customerServices) == 0 {
		return nil, errors.NewError(errors.NotFoundError, nil, "customer service não encontrado")
	}

	return customerServices, nil
}

func (r *CustomerServiceRepository) Update(ctx context.Context, customerService *customerservice.CustomerService) *errors.ErrorOutput {
	defer func() {
		if rec := recover(); rec != nil {
			// Captura e loga a mensagem do panic
			log.Printf("panic occurred: %v", rec)

			// Retorna um erro detalhado com a mensagem do panic
			detailedError := fmt.Sprintf("panic occurred: %v", rec)
			err := errors.NewError(errors.InternalServerError, nil, detailedError)
			log.Printf("Detailed error: %v", err.LogMessage)
		}
	}()

	result, err := r.DB.ExecContext(ctx,
		`UPDATE customer_service SET service_id = $1, customer_id = $2, name = $3, document = $4, email = $5, phone = $6, gender = $7 WHERE id = $8`,
		customerService.ServiceID,
		customerService.CustomerID,
		customerService.Name,
		customerService.Document,
		customerService.Email,
		customerService.Phone,
		customerService.Gender,
		customerService.ID,
	)

	fmt.Println(result)
	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}

func (r *CustomerServiceRepository) Delete(ctx context.Context, customerService *customerservice.CustomerService) *errors.ErrorOutput {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE customer_service SET deleted_at = now() WHERE id = $1`,
		customerService.ID,
	)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}

func (r *CustomerServiceRepository) GetByCustomerAndServiceID(ctx context.Context, customerID, serviceID string) (*customerservice.CustomerService, *errors.ErrorOutput) {
	customerService := &customerservice.CustomerService{}

	err := r.DB.QueryRowContext(ctx,
		`SELECT id, service_id, customer_id, name, document, birth_date, email, phone, gender, created_at from customer_service WHERE customer_id = $1 AND service_id = $2 AND deleted_at IS NULL`,
		customerID,
		serviceID,
	).Scan(
		&customerService.ID,
		&customerService.ServiceID,
		&customerService.CustomerID,
		&customerService.Name,
		&customerService.Document,
		&customerService.BirthDate,
		&customerService.Email,
		&customerService.Phone,
		&customerService.Gender,
		&customerService.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewError(errors.NotFoundError, nil, "customer service não encontrado")
	}

	if err != nil {
		return nil, errors.NewError(errors.InternalServerError, nil, "Error obtendo customer service")
	}

	return customerService, nil
}
