package customerservice_usecase

import (
	"context"
	customerservice "payment-layer-card-api/entities/customer_service"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCustomerServiceByCustomerAndServiceID struct {
	customerServiceRepo customerservice.CustomerServiceRepository
}

func NewGetCustomerServiceByCustomerAndServiceID(customerServiceRepo customerservice.CustomerServiceRepository) *GetCustomerServiceByCustomerAndServiceID {
	return &GetCustomerServiceByCustomerAndServiceID{
		customerServiceRepo: customerServiceRepo,
	}
}

func (g *GetCustomerServiceByCustomerAndServiceID) Execute(ctx context.Context, customerID, serviceID string) (*customerservice.CustomerService, *errors.ErrorOutput) {
	customerService, err := g.customerServiceRepo.GetByCustomerAndServiceID(ctx, customerID, serviceID)
	if err != nil {
		if err.Code == errors.NotFoundError {
			return nil, nil
		}
		return nil, err
	}

	return customerService, nil
}
