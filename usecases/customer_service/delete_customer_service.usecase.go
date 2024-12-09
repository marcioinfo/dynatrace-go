package customerservice_usecase

import (
	"context"
	customerservice "payment-layer-card-api/entities/customer_service"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type DeleteCustomerService struct {
	CustomerServiceRepository customerservice.CustomerServiceRepository
}

func NewDeleteCustomerService(repo customerservice.CustomerServiceRepository) *DeleteCustomerService {
	return &DeleteCustomerService{
		CustomerServiceRepository: repo,
	}
}

func (d *DeleteCustomerService) Execute(ctx context.Context, customerServiceID string) (*customerservice.CustomerService, *errors.ErrorOutput) {
	customerService := &customerservice.CustomerService{
		ID: customerServiceID,
	}

	err := d.CustomerServiceRepository.Delete(ctx, customerService)
	if err != nil {
		return nil, err
	}

	return &customerservice.CustomerService{
		ID: customerServiceID,
	}, nil
}
