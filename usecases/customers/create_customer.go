package customers_usecase

import (
	"context"
	"payment-layer-card-api/common/helpers"
	customer_entity "payment-layer-card-api/entities/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CreateCustomer struct {
	CustomerRepository customer_entity.CustomerRepositoryInterface
}

func NewCreateCustomer(customerRepository customer_entity.CustomerRepositoryInterface) *CreateCustomer {
	return &CreateCustomer{
		CustomerRepository: customerRepository,
	}
}

func (c *CreateCustomer) Execute(ctx context.Context, input customer_dto.CreateCustomerDTOInput) (*customer_dto.CreateCustomerDTOOutput, *errors.ErrorOutput) {
	if input.Document != "" {
		DocumentIsValid := helpers.DocumentIsValid(input.Document)
		if !DocumentIsValid {
			return nil, errors.NewError(errors.DocumentIsInvalid, nil, "Document é inválido.")
		}
	}

	customerAlreadyExists, errGetCustomerByEmail := c.getCustomerByEmail(ctx, input.Email)

	if errGetCustomerByEmail != nil {
		return nil, errGetCustomerByEmail
	}

	if customerAlreadyExists != nil {
		return nil, errors.NewError(errors.CustomerAlreadyExists, nil, "E-mail informado já registrado em nosso sistema. Por favor digite outro!")
	}

	phone := helpers.AddPrefixToNumber(input.Phone)

	customer := customer_entity.NewCustomer()
	customer.InitID()

	customer.Phone = phone
	customer.Name = input.FirstName + " " + input.LastName
	customer.Document = input.Document
	customer.Email = input.Email
	customer.Gender = input.Gender
	customer.BirthDate = input.BirthDate

	valid := customer.IsValid()
	if valid != nil {
		return nil, errors.NewError(errors.ValidationEntityError, valid)
	}

	errorCreateCustomer := c.createCustomer(ctx, customer)
	if errorCreateCustomer != nil {
		return nil, errorCreateCustomer
	}

	return &customer_dto.CreateCustomerDTOOutput{
		ID: customer.ID,
	}, nil
}

func (c *CreateCustomer) createCustomer(ctx context.Context, customer *customer_entity.Customer) *errors.ErrorOutput {
	err := c.CustomerRepository.Insert(ctx, customer)
	if err != nil {
		return err
	}

	return nil
}

func (c *CreateCustomer) getCustomerByEmail(ctx context.Context, email string) (*customer_entity.Customer, *errors.ErrorOutput) {
	customer, err := c.CustomerRepository.GetCustomerByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
