package customers_usecase

import (
	"context"
	"payment-layer-card-api/entities/customers"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCustomerByEmail struct {
	CustomerRepository customers.CustomerRepositoryInterface
}

func NewGetCustomerByEmail(customerRepository customers.CustomerRepositoryInterface) *GetCustomerByEmail {
	return &GetCustomerByEmail{
		CustomerRepository: customerRepository,
	}
}

func (g *GetCustomerByEmail) Execute(ctx context.Context, email string) (*customers.Customer, *errors.ErrorOutput) {
	customer, err := g.CustomerRepository.GetCustomerByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
