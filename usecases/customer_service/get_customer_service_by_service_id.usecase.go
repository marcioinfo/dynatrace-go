package customerservice_usecase

import (
	"context"
	customerservice "payment-layer-card-api/entities/customer_service"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCustomerServiceByServiceID struct {
	CustomerServiceRepository customerservice.CustomerServiceRepository
}

func NewGetCustomerServiceByServiceID(repo customerservice.CustomerServiceRepository) *GetCustomerServiceByServiceID {
	return &GetCustomerServiceByServiceID{
		CustomerServiceRepository: repo,
	}
}

func (g *GetCustomerServiceByServiceID) Execute(ctx context.Context, serviceID string) ([]*customerservice.CustomerService, *errors.ErrorOutput) {
	customerServices, err := g.CustomerServiceRepository.GetByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	return customerServices, nil
}
