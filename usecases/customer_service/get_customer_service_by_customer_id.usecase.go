package customerservice_usecase

import (
	"context"
	customerservice "payment-layer-card-api/entities/customer_service"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCustomerServiceByCustomerID struct {
	CustomerServiceRepository customerservice.CustomerServiceRepository
}

func NewGetCustomerServiceByCustomerID(repo customerservice.CustomerServiceRepository) *GetCustomerServiceByCustomerID {
	return &GetCustomerServiceByCustomerID{
		CustomerServiceRepository: repo,
	}
}

func (g *GetCustomerServiceByCustomerID) Execute(ctx context.Context, customerID string) ([]*customerservice.CustomerService, *errors.ErrorOutput) {
	customerServices, err := g.CustomerServiceRepository.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return customerServices, nil
}
