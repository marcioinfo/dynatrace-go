package http

import (
	"fmt"
	"net/http"
	"os"
	"payment-layer-card-api/bootstrap"
	card_token_usecase "payment-layer-card-api/usecases/card_token"
	cards_usecase "payment-layer-card-api/usecases/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"
	"payment-layer-card-api/usecases/customer_token"
	customers_usecase "payment-layer-card-api/usecases/customers"
	card_workflow "payment-layer-card-api/workflows/card"
	"runtime/debug"

	"payment-layer-card-api/workflows/customer"

	cardvalidator "github.com/adhfoundation/layer-tools/card-validator"
	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/layer-tools/tenant"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CardController struct {
	app *bootstrap.App
}

func NewCardController(app *bootstrap.App) *CardController {
	return &CardController{
		app: app,
	}
}

// TODO: MUDAR PARA SNAKE CASE
func ChooseCreateCardField(createCardInput *card_dtos.CreateCardDTOInput) *card_dtos.CreateCardDTOInput {
	if createCardInput.Customer.Address.PostalCode != "" {
		createCardInput.Customer.Address.PostalCodeAlt = ""
	} else {
		if createCardInput.Customer.Address.PostalCodeAlt != "" {
			createCardInput.Customer.Address.PostalCode = createCardInput.Customer.Address.PostalCodeAlt
			createCardInput.Customer.Address.PostalCodeAlt = ""
		}
	}
	return createCardInput
}

// @Summary Criar um cartão
// @Description Cria um cartão e o tokemiza nos gateways de pagamento
// @Tags Cards
// @Accept  json
// @Produce  json
// @Param card body card_dtos.CreateCardDTOInput true "Informações sobre o cartão e informações do pagador"
// @Success 201 {object} card_dtos.CreateCardDTOOutput
// @Failure 400 {object} ErrorWrapperInternalDTO
// @Router /cards [post]
// @Security x-api-key
func (cc *CardController) CreateCard(c echo.Context) error {
	createCardDTOInput := &card_dtos.CreateCardDTOInput{}
	ctx := c.Request().Context()

	err := c.Bind(createCardDTOInput)
	if err != nil {
		errorController := errors.NewError(errors.ReadBody, err, "Erro desserializando DTO")
		log.Error(ctx, err).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.LogMessage)
	}

	//TODO: MUDAR PARA SNAKE CASE

	if createCardDTOInput.CustomerID == "" && createCardDTOInput.Customer == nil {
		errInstance := errors.NewError(errors.BadRequest, fmt.Errorf("customer_id or customer is required"), "customer_id ou customer é obrigatório")
		log.Error(ctx, errInstance.LogMessageToError()).Msg(errInstance.Message)
		return errors.EchoErrorResponse(c, *errInstance, errInstance.Message)
	}

	var serviceId []string
	if createCardDTOInput.Customer != nil {
		createCardDTOInput = ChooseCreateCardField(createCardDTOInput)
		serviceId = append(serviceId, createCardDTOInput.Customer.ServiceID)
	} else {
		customerService, err := cc.app.CustomerServiceRepo.GetByCustomerID(ctx, createCardDTOInput.CustomerID)
		if err != nil {
			return errors.EchoErrorResponse(c, *err, err.Message)
		}
		for _, service := range customerService {
			serviceId = append(serviceId, service.ServiceID)
		}
	}

	urlIntegration := os.Getenv("INTEGRATION_API_URL")
	key := os.Getenv("INTEGRATION_API_KEY")

	isAllowed, errTenant := tenant.VerifyTenant(ctx, key, urlIntegration, serviceId)
	if errTenant != nil {
		log.Error(ctx, errTenant).Msg("Erro ao verificar tenant: " + errTenant.Error())
		errInstance := errors.NewError(errors.InternalServerError, errTenant, "Erro ao verificar identidade")
		return errors.EchoErrorResponse(c, *errInstance, errInstance.Message)
	}

	if !isAllowed {
		errInstance := errors.NewError(errors.Unauthorized, fmt.Errorf("apikey has not permission to create card with this service"), "Você não tem permissão para criar cartao para este cliente")
		log.Error(ctx, errInstance.LogMessageToError()).Msg(errInstance.LogMessage)
		return errors.EchoErrorResponse(c, *errInstance, errInstance.Message)
	}

	err = cc.app.Validator.Struct(createCardDTOInput)
	if err != nil {
		return errors.ValidationErrorResponse(c, err)
	}

	if createCardDTOInput.Brand == "" {
		result := cardvalidator.CardNumber(createCardDTOInput.Number, cardvalidator.CardNumberOptions{})
		if result.IsValid {
			createCardDTOInput.Brand = result.Card.Type
		}
	}

	customersUseCaseGetCustomerById := customers_usecase.NewGetCustomerByID(cc.app.CustomerRepo)
	customerTokenUseCaseGetByCustomerIDAndGateway := customer_token.NewGetByCustomerIDAndGateway(cc.app.CustomerTokenRepo)
	cardsUseCaseVerifyCardAlreadyExists := cards_usecase.NewVerifyCardAlreadyExists(cc.app.CardRepo)
	cardsUseCaseCreateCard := cards_usecase.NewCreateCard(cc.app.CardRepo)
	cardTokenUseCaseGetByCardIDAndGateway := card_token_usecase.NewGetCardTokenByCardIDAndGateway(cc.app.CardTokenRepo)
	customerWorkflow := customer.NewCreateCustomerWorkflow(cc.app)
	createCardWorkflow := card_workflow.NewCreateCardWorkFlow(
		cc.app.PaymentGateways,
		customersUseCaseGetCustomerById,
		customerTokenUseCaseGetByCustomerIDAndGateway,
		cardsUseCaseVerifyCardAlreadyExists,
		cardsUseCaseCreateCard,
		cardTokenUseCaseGetByCardIDAndGateway,
		cc.app.CardTokenRepo,
		customerWorkflow,
		cc.app.GetCardRetryCount(),
		cc.app.GetCardRetryDelayInMilliseconds(),
		cc.app.QueueService,
		cc.app.GetSyncCardQueueURL(),
	)

	createCardDTOOutput, errorCreateCardWorkflow := createCardWorkflow.Execute(ctx, createCardDTOInput)
	if errorCreateCardWorkflow != nil {
		log.Info(c.Request().Context()).Msg("Ocorreu um erro ao criar o cartão!")
		return errors.EchoErrorResponse(c, *errorCreateCardWorkflow, errorCreateCardWorkflow.Message)
	}

	log.Info(ctx).Msg("Cartão criado com sucesso!")
	return c.JSON(http.StatusCreated, createCardDTOOutput)
}

// @Summary Obter um cartão pelo id
// @Description Buscar um cartão pelo identificador
// @Tags Cards
// @Accept  json
// @Produce  json
// @Param id path string true "Identificador do cartão"
// @Success 200 {object} card_dtos.GetCardWithTokensDTO
// @Failure 400 {object} ErrorWrapperInternalDTO
// @Router /cards/{id} [get]
// @Security x-api-key
func (cc *CardController) GetCardByID(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			ctx := c.Request().Context()
			errMessage := "Recuperação de pânico: "
			if err, ok := r.(error); ok {
				errMessage += err.Error()
			} else {
				errMessage += "Erro desconhecido"
			}
			errMessage += "\nStack trace:\n" + string(debug.Stack())
			err := errors.NewError(errors.InternalServerError, nil, errMessage)
			log.Error(ctx, err).Msg("Error ocorreu no GetCardByID")
			c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Um erro interno aconteceu. Por favor tente novamente mais tarde.",
			})
		}
	}()

	id := c.Param("id")
	ctx := c.Request().Context()
	uuidErr := uuid.Validate(id)
	if id == "" || uuidErr != nil {
		errorController := errors.NewError(errors.ParameterIsRequired, uuidErr, "ID campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	getCardByIdWorkflow := card_workflow.NewGetCardByIdWorkflow(cc.app)
	getCardWithTokensDTO, errorGetCardByIdWorkflow := getCardByIdWorkflow.Execute(ctx, id)
	if errorGetCardByIdWorkflow != nil {
		return errors.EchoErrorResponse(c, *errorGetCardByIdWorkflow, errorGetCardByIdWorkflow.Message)
	}

	log.Info(ctx).Msg("Cartão obtido com sucesso pelo ID!")
	return c.JSON(http.StatusOK, getCardWithTokensDTO)
}

// @Summary Obter um cartão pelo id e id do cliente
// @Description Buscar um cartão pelo identificador e id do cliente
// @Tags Cards
// @Accept  json
// @Produce  json
// @Param cardId path string true "Identificador do cartão"
// @Param customerId path string true "Identificador do cliente"
// @Success 200 {object} card_dtos.GetCardWithoutTokensDTO
// @Failure 400 {object} ErrorWrapperInternalDTO
// @Router /cards/{cardId}/customer/{customerId} [get]
// @Security x-api-key
func (cc *CardController) GetCardByIDAndCustomerID(c echo.Context) error {
	ctx := c.Request().Context()

	cardId := c.Param("cardId")
	uuidErr := uuid.Validate(cardId)
	if cardId == "" || uuidErr != nil {
		errorController := errors.NewError(errors.ParameterIsRequired, uuidErr, "CardID campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	customerId := c.Param("customerId")
	uuidErr = uuid.Validate(customerId)
	if customerId == "" || uuidErr != nil {
		errorController := errors.NewError(errors.ParameterIsRequired, uuidErr, "CustomerID campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	getCardByIdAndCustomerIDUsecase := cards_usecase.NewGetCardByIDAndCustomerID(cc.app.CardRepo)

	card, err := getCardByIdAndCustomerIDUsecase.Execute(ctx, cardId, customerId)
	if err != nil {
		return errors.EchoErrorResponse(c, *err, err.Message)
	}

	cardDto := &card_dtos.GetCardWithoutTokensDTO{
		ID:          card.ID,
		CustomerID:  card.CustomerID,
		Holder:      card.Holder,
		Brand:       card.Brand,
		FirstDigits: card.FirstDigits,
		LastDigits:  card.LastDigits,
		ExpMonth:    card.ExpMonth,
		ExpYear:     card.ExpYear,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}

	log.Info(ctx).Msg("Cartão obtido com sucesso pelo ID!")
	return c.JSON(http.StatusOK, cardDto)
}

func (cc *CardController) DeleteCardByIDAndCustomerID(c echo.Context) error {
	ctx := c.Request().Context()

	cardId := c.Param("cardId")
	uuidErr := uuid.Validate(cardId)
	if cardId == "" || uuidErr != nil {
		errorController := errors.NewError(errors.ParameterIsRequired, uuidErr, "CardID campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	customerId := c.Param("customerId")
	uuidErr = uuid.Validate(customerId)
	if customerId == "" || uuidErr != nil {
		errorController := errors.NewError(errors.ParameterIsRequired, uuidErr, "CustomerID campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	var serviceId []string

	customerService, err := cc.app.CustomerServiceRepo.GetByCustomerID(ctx, customerId)
	if err != nil {
		return errors.EchoErrorResponse(c, *err, err.Message)
	}
	for _, service := range customerService {
		serviceId = append(serviceId, service.ServiceID)
	}

	urlIntegration := os.Getenv("INTEGRATION_API_URL")
	key := os.Getenv("INTEGRATION_API_KEY")

	isAllowed, errTenant := tenant.VerifyTenant(ctx, key, urlIntegration, serviceId)
	if errTenant != nil {
		log.Error(ctx, errTenant).Msg("Erro ao verificar tenant: " + errTenant.Error())
		errInstance := errors.NewError(errors.InternalServerError, errTenant, "Erro ao verificar identidade")
		return errors.EchoErrorResponse(c, *errInstance, errInstance.Message)
	}

	if !isAllowed {
		errInstance := errors.NewError(errors.Unauthorized, fmt.Errorf("apikey has not permission to create card with this service"), "Você não tem permissão para deletar esse cartão")
		log.Error(ctx, errInstance.LogMessageToError()).Msg(errInstance.LogMessage)
		return errors.EchoErrorResponse(c, *errInstance, errInstance.Message)
	}

	customerTokenUseCaseGetByCustomerIDAndGateway := customer_token.NewGetByCustomerIDAndGateway(cc.app.CustomerTokenRepo)
	deleteCardUsecase := cards_usecase.NewDeleteCardByID(cc.app.CardRepo)
	getCardByIdAndCustomerIdUsecase := cards_usecase.NewGetCardByIDAndCustomerID(cc.app.CardRepo)
	cardTokenUseCaseGetByCardIDAndGateway := card_token_usecase.NewGetCardTokenByCardIDAndGateway(cc.app.CardTokenRepo)

	deleteCardWorkflow := card_workflow.NewDeleteCardWorkflow(
		cc.app.PaymentGateways,
		getCardByIdAndCustomerIdUsecase,
		cardTokenUseCaseGetByCardIDAndGateway,
		customerTokenUseCaseGetByCustomerIDAndGateway,
		deleteCardUsecase,
		cc.app.CardTokenRepo,
	)

	response := deleteCardWorkflow.Execute(ctx, cardId, customerId)

	if response.Error == nil {
		log.Info(ctx).Msg("Cartão deletado com sucesso pelo ID!")
	} else {
		log.Info(ctx).Msgf("Erro ao tentar deletar card por ID: %+v", response.Error)
	}

	if len(response.TokensDeleted) == 0 {
		return errors.EchoErrorResponse(c, *response.Error, "Error ao tentar deletar cartão")
	}

	return c.JSON(http.StatusOK, response)

}
