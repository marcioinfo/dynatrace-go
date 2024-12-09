package customers_usecase

import (
	"context"
	"payment-layer-card-api/common/helpers"
	"payment-layer-card-api/entities/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type UpdateCustomerUsecase struct {
	customerRepository customers.CustomerRepositoryInterface
}

func NewUpdateCustomerUsecase(customerRepository customers.CustomerRepositoryInterface) *UpdateCustomerUsecase {
	return &UpdateCustomerUsecase{
		customerRepository: customerRepository,
	}
}

func (uc *UpdateCustomerUsecase) Execute(ctx context.Context, input *customer_dto.UpdateCustomerDTO, id string) (*customers.Customer, *errors.ErrorOutput) {
	existingCustomer, err := uc.customerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Email != nil {
		existsEmail, _ := uc.customerRepository.GetCustomerByEmail(ctx, *input.Email)
		if existsEmail != nil && existsEmail.ID != id {
			return nil, errors.NewError(errors.AlredyExists, nil, "E-mail já cadastrado.")
		}
	}

	if input.Name != nil {
		existingCustomer.Name = *input.Name
	}
	if input.Email != nil {
		existingCustomer.Email = *input.Email
	}
	if input.Phone != nil {
		phone := helpers.AddPrefixToNumber(*input.Phone)
		existingCustomer.Phone = phone
	}
	if input.Gender != nil {
		existingCustomer.Gender = *input.Gender
	}
	if input.BirthDate != nil {
		existingCustomer.BirthDate = *input.BirthDate
	}

	err = uc.customerRepository.Update(ctx, existingCustomer)
	if err != nil {
		return nil, err
	}
	return existingCustomer, nil
}
