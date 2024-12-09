package customerservice_usecase

import (
	"context"
	customerservice "payment-layer-card-api/entities/customer_service"
	"time"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type UpdateCustomerService struct {
	CustomerServiceRepository customerservice.CustomerServiceRepository
}

func NewUpdateCustomerService(repo customerservice.CustomerServiceRepository) *UpdateCustomerService {
	return &UpdateCustomerService{
		CustomerServiceRepository: repo,
	}
}

func (u *UpdateCustomerService) Execute(ctx context.Context, input *customerservice.CustomerService) (*customerservice.CustomerService, *errors.ErrorOutput) {
	input.UpdatedAt = datetypes.CustomDateTime(time.Now())

	if input.ServiceID == "" {
		return nil, errors.NewError(errors.ParameterIsRequired, nil, "Service ID parâmetro é obrigatório.")
	}

	if input.CustomerID == "" {
		return nil, errors.NewError(errors.ParameterIsRequired, nil, "Customer ID parâmetro é obrigatório.")
	}

	err := u.CustomerServiceRepository.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return &customerservice.CustomerService{
		ID: input.ID,
	}, nil
}
