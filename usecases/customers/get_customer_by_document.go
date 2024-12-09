package customers_usecase

import (
	"context"
	"payment-layer-card-api/common/helpers"
	"payment-layer-card-api/entities/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCustomerByDocument struct {
	CustomerRepository customers.CustomerRepositoryInterface
}

func NewGetCustomerByDocument(customerRepository customers.CustomerRepositoryInterface) *GetCustomerByDocument {
	return &GetCustomerByDocument{
		CustomerRepository: customerRepository,
	}
}

func (g *GetCustomerByDocument) Execute(ctx context.Context, document string, serviceID string) (*customer_dto.CustomerOutputDTO, *layerErrors.ErrorOutput) {
	if document == "" || !helpers.DocumentIsValid(document) {
		return nil, layerErrors.NewError(layerErrors.ParameterIsRequired, nil, "Document parâmetro é inválido.")
	}

	customer, err := g.getCustomer(ctx, document, serviceID)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (g *GetCustomerByDocument) getCustomer(ctx context.Context, document string, serviceID string) (*customer_dto.CustomerOutputDTO, *layerErrors.ErrorOutput) {
	customer, err := g.CustomerRepository.GetCustomerByDocument(ctx, document, serviceID)

	if err != nil {
		return nil, err
	}

	return &customer_dto.CustomerOutputDTO{
		ID:        customer.ID,
		Name:      customer.Name,
		Document:  customer.Document,
		Email:     customer.Email,
		Phone:     customer.Phone,
		BirthDate: customer.BirthDate,
		Gender:    customer.Gender,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}, nil
}
