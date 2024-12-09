package customer

import (
	"context"
	"payment-layer-card-api/bootstrap"
	customertoken_usecase "payment-layer-card-api/usecases/customer_token"
	customers_usecase "payment-layer-card-api/usecases/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCustomerByIDWorkflow struct {
	app *bootstrap.App
}

func NewGetCustomerByIdWorkflow(app *bootstrap.App) *GetCustomerByIDWorkflow {
	return &GetCustomerByIDWorkflow{
		app: app,
	}
}

func (gcw *GetCustomerByIDWorkflow) Execute(ctx context.Context, id string) (*customer_dto.CustomerWithTokensOutputDTO, *errors.ErrorOutput) {
	getCustomerByIDUsecase := customers_usecase.NewGetCustomerByID(gcw.app.CustomerRepo)

	customer, err := getCustomerByIDUsecase.Execute(ctx, id)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, errors.NewPaymentLayerError(errors.CustomerNotFound, "Customer não encontrado.")
	}

	customerTokens, err := gcw.getCustomerTokens(ctx, customer.ID)
	if err != nil {
		return nil, err
	}

	return &customer_dto.CustomerWithTokensOutputDTO{
		ID:        customer.ID,
		Name:      customer.Name,
		Document:  customer.Document,
		Email:     customer.Email,
		Phone:     customer.Phone,
		BirthDate: customer.BirthDate,
		Gender:    customer.Gender,
		Tokens:    customerTokens,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}, nil

}

func (gcw *GetCustomerByIDWorkflow) getCustomerTokens(ctx context.Context, customerID string) ([]*customer_dto.CustomerTokens, *errors.ErrorOutput) {
	getCustomerTokensUsecase := customertoken_usecase.NewGetByCustomerID(gcw.app.CustomerTokenRepo)

	customerTokens, err := getCustomerTokensUsecase.Execute(ctx, customerID)
	if err != nil {
		return nil, err
	}

	var customerTokensOutput []*customer_dto.CustomerTokens
	for _, customerTokenIt := range customerTokens {
		customerToken := &customer_dto.CustomerTokens{}

		customerToken.Gateway = customerTokenIt.Gateway
		customerToken.Token = customerTokenIt.CustomerToken

		customerTokensOutput = append(customerTokensOutput, customerToken)
	}

	return customerTokensOutput, nil
}
