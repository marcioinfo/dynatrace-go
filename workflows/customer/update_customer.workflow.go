package customer

import (
	"context"
	"encoding/json"
	"os"
	"payment-layer-card-api/bootstrap"
	customer_token_entity "payment-layer-card-api/entities/customer_token"
	"payment-layer-card-api/entities/customers"
	"payment-layer-card-api/entities/payment_gateway"
	"payment-layer-card-api/usecases/customer_token"
	customers_usecase "payment-layer-card-api/usecases/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"go.elastic.co/apm/v2"
)

type UpdateCustomerWorkFlow struct {
	app *bootstrap.App
}

func NewUpdateCustomerWorkFlow(app *bootstrap.App) *UpdateCustomerWorkFlow {
	return &UpdateCustomerWorkFlow{
		app: app,
	}
}

func (ucw *UpdateCustomerWorkFlow) Execute(ctx context.Context, input *customer_dto.UpdateCustomerDTO, id string) (*customers.Customer, *errors.ErrorOutput) {
	customer, err := ucw.updateCustomer(ctx, input, id)
	if err != nil {
		return nil, err
	}
	err = ucw.updateCustomerTokens(ctx, customer, id)
	if err != nil {
		return nil, err
	}
	err = ucw.sendMessage(ctx, customer)
	if err != nil {

		return nil, err
	}

	return customer, nil
}

func (ucw *UpdateCustomerWorkFlow) updateCustomer(ctx context.Context, input *customer_dto.UpdateCustomerDTO, id string) (*customers.Customer, *errors.ErrorOutput) {
	updateCustomerUsecase := customers_usecase.NewUpdateCustomerUsecase(ucw.app.CustomerRepo)
	return updateCustomerUsecase.Execute(ctx, input, id)
}

func (ucw *UpdateCustomerWorkFlow) updateCustomerTokens(ctx context.Context, customer *customers.Customer, id string) *errors.ErrorOutput {
	for _, gateway := range ucw.app.PaymentGateways {
		customerTokens, err := ucw.getCustomerTokens(ctx, id)
		if err != nil {
			return err
		}

		err = ucw.updateTokensForGateway(ctx, customerTokens, gateway, customer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ucw *UpdateCustomerWorkFlow) getCustomerTokens(ctx context.Context, customerID string) ([]*customer_token_entity.CustomerToken, *errors.ErrorOutput) {
	usecaseToken := customer_token.NewGetByCustomerID(ucw.app.CustomerTokenRepo)
	return usecaseToken.Execute(ctx, customerID)
}

func (ucw *UpdateCustomerWorkFlow) updateTokensForGateway(ctx context.Context, customerTokens []*customer_token_entity.CustomerToken, gateway payment_gateway.PaymentGatewayInterface, customer *customers.Customer) *errors.ErrorOutput {
	for _, customerToken := range customerTokens {
		if customerToken.Gateway == gateway.GatewayName() {
			_, _, err := gateway.UpdateCustomer(ctx, customer, customerToken.CustomerToken)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ucw *UpdateCustomerWorkFlow) sendMessage(ctx context.Context, customer *customers.Customer) *errors.ErrorOutput {
	tx := apm.TransactionFromContext(ctx)
	if tx != nil {
		customer.ApmLink.SpanId = tx.TraceContext().Span
		customer.ApmLink.TraceId = tx.TraceContext().Trace
	}

	updateCustomerSQSPayload, err := json.Marshal(customer)
	if err != nil {
		return errors.NewPaymentLayerError(errors.InternalServerError, err.Error())
	}

	queueUrl := os.Getenv("UPDATE_CUSTOMER_INTEGRATION_QUEUE_URL")

	err = ucw.app.QueueService.SendMessageWithContext(ctx, queueUrl, string(updateCustomerSQSPayload))
	if err != nil {
		return errors.NewPaymentLayerError(errors.InternalServerError, err.Error())
	}
	return nil
}
