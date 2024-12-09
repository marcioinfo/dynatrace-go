package customerservice_usecase

import (
	"context"
	"payment-layer-card-api/common/helpers"
	customerservice "payment-layer-card-api/entities/customer_service"
	"time"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CreateCustomerService struct {
	CustomerServiceRepository customerservice.CustomerServiceRepository
}

func NewCreateCustomerService(repo customerservice.CustomerServiceRepository) *CreateCustomerService {
	return &CreateCustomerService{
		CustomerServiceRepository: repo,
	}
}

func (c *CreateCustomerService) Execute(ctx context.Context, input *customerservice.CustomerService) (*customerservice.CustomerService, *errors.ErrorOutput) {
	if input.Document != "" {
		DocumentIsValid := helpers.DocumentIsValid(input.Document)
		if !DocumentIsValid {
			return nil, errors.NewError(errors.DocumentIsInvalid, nil, "Document é inválido.")
		}
	}

	input.InitID()
	input.UpdatedAt = datetypes.CustomDateTime(time.Now())

	errorCreateCustomerService := c.createCustomerService(ctx, input)
	if errorCreateCustomerService != nil {
		return nil, errorCreateCustomerService
	}

	return &customerservice.CustomerService{
		ID: input.ID,
	}, nil
}

func (c *CreateCustomerService) createCustomerService(ctx context.Context, customerService *customerservice.CustomerService) *errors.ErrorOutput {
	err := c.CustomerServiceRepository.Insert(ctx, customerService)
	if err != nil {
		return err
	}

	return nil
}
