package customer

import (
	"context"
	"encoding/json"
	"os"
	"payment-layer-card-api/bootstrap"
	customerservice "payment-layer-card-api/entities/customer_service"

	customerservice_usecase "payment-layer-card-api/usecases/customer_service"
	customer_token_usecase "payment-layer-card-api/usecases/customer_token"
	customer_usecase "payment-layer-card-api/usecases/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"go.elastic.co/apm/v2"
)

type CreateCustomerWorkflow struct {
	app *bootstrap.App
}

func NewCreateCustomerWorkflow(app *bootstrap.App) *CreateCustomerWorkflow {
	return &CreateCustomerWorkflow{
		app: app,
	}
}

func (c *CreateCustomerWorkflow) Execute(ctx context.Context, input *customer_dto.CreateCustomerDTOInput) (*customer_dto.CreateCustomerDTOOutput, *errors.ErrorOutput) {
	customerDTOOutput, errorFindOrCreateCustomer := c.findOrCreateCustomer(ctx, input)
	if errorFindOrCreateCustomer != nil {
		return nil, errorFindOrCreateCustomer
	}

	errorCreateCustomerTokens := c.createCustomerTokens(ctx, customerDTOOutput.ID, input)
	if errorCreateCustomerTokens != nil {
		return nil, errorCreateCustomerTokens
	}

	errorSendMessage := c.sendMessage(ctx, customerDTOOutput.ID, input.ServiceID)
	if errorSendMessage != nil {
		return nil, errorSendMessage
	}

	return customerDTOOutput, nil
}

func (c *CreateCustomerWorkflow) findOrCreateCustomer(ctx context.Context, input *customer_dto.CreateCustomerDTOInput) (*customer_dto.CreateCustomerDTOOutput, *errors.ErrorOutput) {
	getCustomerByEmailUsecase := customer_usecase.NewGetCustomerByEmail(c.app.CustomerRepo)
	getCustomerService := customerservice_usecase.NewGetCustomerServiceByCustomerAndServiceID(c.app.CustomerServiceRepo)
	createCustomerServiceUsecase := customerservice_usecase.NewCreateCustomerService(c.app.CustomerServiceRepo)
	updateCustomerUseCase := customer_usecase.NewUpdateCustomerUsecase(c.app.CustomerRepo)
	updateCustomerService := customerservice_usecase.NewUpdateCustomerService(c.app.CustomerServiceRepo)

	customerExists, err := getCustomerByEmailUsecase.Execute(ctx, input.Email)
	if err != nil && err.Code == errors.InternalServerError {
		log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao buscar o Customer pelo documento (%v). Erro (%v)", input.Document, err)
		return nil, err
	}

	if customerExists != nil {
		customerServiceInput := &customerservice.CustomerService{
			ServiceID:  input.ServiceID,
			CustomerID: customerExists.ID,
			Name:       input.FirstName + " " + input.LastName,
			Document:   input.Document,
			BirthDate:  input.BirthDate,
			Email:      input.Email,
			Phone:      input.Phone,
			Gender:     input.Gender,
		}
		customerServiceExisted, err := getCustomerService.Execute(ctx, customerExists.ID, input.ServiceID)
		if err != nil {
			log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao buscar o CustomerService pelo CustomerID (%v) e ServiceID (%v). Erro (%v)", customerExists.ID, input.ServiceID, err)
		}

		if customerServiceExisted != nil {
			customerServiceInput.ID = customerServiceExisted.ID
			_, err = updateCustomerService.Execute(ctx, customerServiceInput)
			if err != nil {
				log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao atualizar o CustomerService pelo CustomerID (%v) e ServiceID (%v). Erro (%v)", customerExists.ID, input.ServiceID, err)
			}

			log.Info(ctx).Msgf("Foi encontrado o CustomerService para o CustomerID (%v) e ServiceID (%v)", customerExists.ID, input.ServiceID)
		} else {
			_, err = createCustomerServiceUsecase.Execute(ctx, customerServiceInput)
			if err != nil {
				log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao criar o CustomerService pelo CustomerID (%v) e ServiceID (%v). Erro (%v)", customerExists.ID, input.ServiceID, err)
			}
		}

		name := input.FirstName + " " + input.LastName
		updateCustomerInput := &customer_dto.UpdateCustomerDTO{
			Name:      &name,
			Email:     &input.Email,
			Phone:     &input.Phone,
			Gender:    &input.Gender,
			BirthDate: &input.BirthDate,
		}

		_, err = updateCustomerUseCase.Execute(ctx, updateCustomerInput, customerExists.ID)
		if err != nil {
			log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao atualizar o Customer pelo documento (%v). Erro (%v)", input.Document, err)
			return nil, err
		}

		log.Info(ctx).Msgf("Foi encontrado o CustomerID (%v) para o documento (%v)", customerExists.ID, input.Document)
		return &customer_dto.CreateCustomerDTOOutput{
			ID: customerExists.ID,
		}, nil
	}

	createCustomerUseCase := customer_usecase.NewCreateCustomer(c.app.CustomerRepo)

	createCustomer, err := createCustomerUseCase.Execute(ctx, *input)
	if err != nil {
		log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao criar o Customer do documento (%v). Erro (%v)", input.Document, err)
		return nil, err
	}

	customerServiceInput := &customerservice.CustomerService{
		ServiceID:  input.ServiceID,
		CustomerID: createCustomer.ID,
		Name:       input.FirstName + " " + input.LastName,
		Document:   input.Document,
		BirthDate:  input.BirthDate,
		Email:      input.Email,
		Phone:      input.Phone,
		Gender:     input.Gender,
	}

	_, err = createCustomerServiceUsecase.Execute(ctx, customerServiceInput)
	if err != nil {
		log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao criar o CustomerService pelo  e ServiceID (%v). Erro (%v)", input.ServiceID, err)
	}

	log.Info(ctx).Msgf("O CustomerID (%v) foi criado para o documento (%v)", createCustomer.ID, input.Document)
	return createCustomer, nil
}

func (c *CreateCustomerWorkflow) createCustomerTokens(ctx context.Context, customerID string, input *customer_dto.CreateCustomerDTOInput) *errors.ErrorOutput {

	countErrors := 0
	var lastError *errors.ErrorOutput
	for _, gateway := range c.app.PaymentGateways {

		gatewayName := gateway.GatewayName()

		createCustomerTokenUseCase := customer_token_usecase.NewCreateCustomerToken(c.app.CustomerTokenRepo, gateway, c.app.GetCustomerRetryCount(), c.app.GetCustomerRetryDelayInMilliseconds())
		_, errorCreateCustomerToken := createCustomerTokenUseCase.Execute(ctx, input, customerID)
		if errorCreateCustomerToken != nil {
			log.Error(ctx, errorCreateCustomerToken.LogMessageToError()).Msgf("Erro ao criar CustomerToken para o CostumerID (%v) no Gateway (%v), retornou o erro (%v)", customerID, gatewayName, errorCreateCustomerToken.Error())
			countErrors++
			lastError = errorCreateCustomerToken
		}
	}

	if countErrors >= len(c.app.PaymentGateways) {
		log.Info(ctx).Msgf("Ocorreu erro na busca ou criação de todos os CustomerTokens para o CustomerId (%v)", customerID)
		return lastError
	}
	return nil
}

func (c *CreateCustomerWorkflow) sendMessage(ctx context.Context, id string, serviceID string) *errors.ErrorOutput {
	customerTokenUsecase := customer_token_usecase.NewGetByCustomerID(c.app.CustomerTokenRepo)
	customerUsecase := customer_usecase.NewGetCustomerByID(c.app.CustomerRepo)
	getCustomerService := customerservice_usecase.NewGetCustomerServiceByCustomerAndServiceID(c.app.CustomerServiceRepo)
	customer, err := customerUsecase.Execute(ctx, id)
	if err != nil {
		return errors.NewPaymentLayerError(errors.InternalServerError, err.Error())
	}

	tokens, err := customerTokenUsecase.Execute(ctx, id)
	if err != nil {
		return errors.NewPaymentLayerError(errors.InternalServerError, err.Error())
	}

	var tokensDTO []customer_dto.Token
	for _, token := range tokens {
		tokensDTO = append(tokensDTO, customer_dto.Token{
			Token:   token.CustomerToken,
			Gateway: token.Gateway,
		})
	}

	customerService, err := getCustomerService.Execute(ctx, customer.ID, serviceID)
	if err != nil {
		return errors.NewPaymentLayerError(errors.InternalServerError, "Erro ao buscar customer service para enviar para fila: "+err.Error())
	}

	payload := &customer_dto.SendCustomerDTO{
		ID:                customer.ID,
		Name:              customer.Name,
		Document:          customer.Document,
		Birthdate:         customer.BirthDate,
		Email:             customer.Email,
		Phone:             customer.Phone,
		Gender:            customer.Gender,
		ServiceID:         serviceID,
		CustomerServiceID: customerService.ID,
		Tokens:            tokensDTO,
	}

	tx := apm.TransactionFromContext(ctx)
	if tx != nil {
		payload.ApmLink.SpanId = tx.TraceContext().Span
		payload.ApmLink.TraceId = tx.TraceContext().Trace
	}

	createCustomerSQSPayload, errResp := json.Marshal(payload)
	if errResp != nil {
		return errors.NewPaymentLayerError(errors.InternalServerError, errResp.Error())
	}

	queueUrl := os.Getenv("CREATE_CUSTOMER_INTEGRATION_QUEUE_URL")

	errResp = c.app.QueueService.SendMessageWithContext(ctx, queueUrl, string(createCustomerSQSPayload))
	if errResp != nil {
		return errors.NewPaymentLayerError(errors.InternalServerError, errResp.Error())
	}
	return nil
}
