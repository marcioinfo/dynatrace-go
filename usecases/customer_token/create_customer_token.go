package customer_token

import (
	"context"
	customer_token_entity "payment-layer-card-api/entities/customer_token"
	"payment-layer-card-api/entities/customers"
	paymentgateway "payment-layer-card-api/entities/payment_gateway"
	customertoken_dto "payment-layer-card-api/usecases/customer_token/dtos"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"strconv"
	"time"

	"github.com/adhfoundation/layer-tools/log"
	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CreateCustomerToken struct {
	CustomerTokenRepository          customer_token_entity.CustomerTokenRepositoryInterface
	PaymentInterface                 paymentgateway.PaymentGatewayInterface
	CustomerRetryCount               int
	CustomerRetryDelayInMilliseconds int
}

func NewCreateCustomerToken(
	customerTokenRepository customer_token_entity.CustomerTokenRepositoryInterface,
	paymentInterface paymentgateway.PaymentGatewayInterface,
	customerRetryCount int,
	customerRetryDelayInMilliseconds int) *CreateCustomerToken {
	return &CreateCustomerToken{
		CustomerTokenRepository:          customerTokenRepository,
		PaymentInterface:                 paymentInterface,
		CustomerRetryCount:               customerRetryCount,
		CustomerRetryDelayInMilliseconds: customerRetryDelayInMilliseconds,
	}
}

func (c *CreateCustomerToken) Execute(ctx context.Context, input *customer_dto.CreateCustomerDTOInput, customer_id string) (*customertoken_dto.CustomerTokenOutputDTO, *layerErrors.ErrorOutput) {
	customerToken := customer_token_entity.NewCustomerToken()
	customerToken.InitID()

	customerTokenExists, err := c.CustomerTokenRepository.GetByCustomerIDAndGateway(ctx, customer_id, c.PaymentInterface.GatewayName())
	if err != nil && err.Code != layerErrors.NotFoundError {
		return nil, err
	}

	if customerTokenExists != nil {
		customer := &customers.Customer{
			ID:        customer_id,
			Name:      input.FirstName + " " + input.LastName,
			Email:     input.Email,
			Document:  input.Document,
			BirthDate: input.BirthDate,
			Phone:     input.Phone,
			Gender:    input.Gender,
		}
		log.Info(ctx).Msgf("O CustomerToken já existe para o CustomerID (%v) no Gateway (%v)", customer_id, c.PaymentInterface.GatewayName())
		createCustomerTokenInGateway, err := c.updateCustomerTokenInGateway(ctx, customer, customerTokenExists.CustomerToken)
		if err != nil {
			return nil, err
		}
		if createCustomerTokenInGateway.Token != customerTokenExists.CustomerToken {

			customerTokenExists.CustomerToken = createCustomerTokenInGateway.Token
			err = c.insertCustomerTokenInDatabase(ctx, customerTokenExists)
			if err != nil {
				return nil, err
			}
		}
		return createCustomerTokenInGateway, nil
	} else {

	}

	createCustomerTokenInGateway, err := c.createCustomerTokenInGateway(ctx, *input, customerToken.ID)
	if err != nil {
		return nil, err
	}

	customerToken.Gateway = createCustomerTokenInGateway.Gateway
	customerToken.CustomerToken = createCustomerTokenInGateway.Token
	customerToken.CustomerId = customer_id

	err = c.insertCustomerTokenInDatabase(ctx, customerToken)
	if err != nil {
		return nil, err
	}

	return createCustomerTokenInGateway, nil
}

func (c *CreateCustomerToken) createCustomerTokenInGateway(ctx context.Context, input customer_dto.CreateCustomerDTOInput, CustomerID string) (*customertoken_dto.CustomerTokenOutputDTO, *layerErrors.ErrorOutput) {
	var customerTokenOutputDTO *customertoken_dto.CustomerTokenOutputDTO
	var httpStatusCode int
	var errorOnCreateCustomer *layerErrors.ErrorOutput

	customerRetryCount := c.CustomerRetryCount
	for i := 1; i <= customerRetryCount; i++ {
		if i > 1 {
			customerRetryDelay := (time.Duration(c.CustomerRetryDelayInMilliseconds) * time.Millisecond)
			log.Info(ctx).Msgf("Aguardando (%v) para a Tentativa (%v) da criação do Token para o CustomerID (%v) no Gateway (%v)", customerRetryDelay, i, CustomerID, c.PaymentInterface.GatewayName())
			time.Sleep(customerRetryDelay)
		}
		log.Info(ctx).Msgf("Tentativa (%v) da criação do Token para o CustomerID (%v) no Gateway (%v)", i, CustomerID, c.PaymentInterface.GatewayName())
		customerTokenOutputDTO, httpStatusCode, errorOnCreateCustomer = c.PaymentInterface.CreateCustomer(ctx, &input, CustomerID)

		if errorOnCreateCustomer == nil {
			log.Info(ctx).Msgf("Token criado com sucesso para o CustomerID (%v) no Gateway (%v) na tentativa (%v)", CustomerID, c.PaymentInterface.GatewayName(), i)
			return customerTokenOutputDTO, nil
		}

		strHttpStatusCode := strconv.Itoa(httpStatusCode)

		log.Error(ctx, errorOnCreateCustomer).Msgf("Ocorreu um erro na criação do Token para o CustomerID (%v) no Gateway (%v) na tentativa (%v) e retornou httpStatus (%v) com o erro (%v) \n", CustomerID, c.PaymentInterface.GatewayName(), i, strHttpStatusCode, errorOnCreateCustomer)
		if strHttpStatusCode[0] == '4' {
			return nil, layerErrors.NewError(layerErrors.InvalidCustomer, nil, "Erro desconhecido ocorreu ao criar customer token")
		}
	}

	log.Error(ctx, errorOnCreateCustomer).Msgf("Não foi possivel criar o token para o CustomerID (%v) no Gateway (%v) depois de (%v) tentativas", CustomerID, c.PaymentInterface.GatewayName(), customerRetryCount)
	return nil, layerErrors.NewError(layerErrors.InvalidCustomer, nil, "Erro desconhecido ocorreu ao criar customer token")
}

func (c *CreateCustomerToken) updateCustomerTokenInGateway(ctx context.Context, customer *customers.Customer, token string) (*customertoken_dto.CustomerTokenOutputDTO, *layerErrors.ErrorOutput) {
	customerTokenOutputDTO, httpStatusCode, errorOnUpdateCustomer := c.PaymentInterface.UpdateCustomer(ctx, customer, token)
	if errorOnUpdateCustomer == nil {
		log.Info(ctx).Msgf("Token atualizado com sucesso para o CustomerID (%v) no Gateway (%v)", customer.ID, c.PaymentInterface.GatewayName())
		return customerTokenOutputDTO, nil
	}

	strHttpStatusCode := strconv.Itoa(httpStatusCode)

	log.Error(ctx, errorOnUpdateCustomer).Msgf("Ocorreu um erro na atualização do Token para o CustomerID (%v) no Gateway (%v) e retornou httpStatus (%v) com o erro (%v) \n", customer.ID, c.PaymentInterface.GatewayName(), strHttpStatusCode, errorOnUpdateCustomer)
	if strHttpStatusCode[0] == '4' {
		return nil, errorOnUpdateCustomer
	}

	log.Error(ctx, errorOnUpdateCustomer).Msgf("Não foi possivel atualizar o token para o CustomerID (%v) no Gateway (%v)", customer.ID, c.PaymentInterface.GatewayName())
	return nil, errorOnUpdateCustomer
}

func (c *CreateCustomerToken) insertCustomerTokenInDatabase(ctx context.Context, customerToken *customer_token_entity.CustomerToken) *layerErrors.ErrorOutput {
	err := c.CustomerTokenRepository.Insert(ctx, customerToken)
	if err != nil {
		log.Error(ctx, err.LogMessageToError()).Msgf("Não foi possivel criar o CustomerToken no Banco de Dados para o CustomerID (%v) no Gateway (%v) ", customerToken.CustomerId, customerToken.Gateway)
		return err
	}

	log.Info(ctx).Msgf("O CustomerToken (%v) foi criado no Banco de Dados para o CustomerID (%v) no Gateway (%v) ", customerToken.ID, customerToken.CustomerId, customerToken.Gateway)
	return nil
}
