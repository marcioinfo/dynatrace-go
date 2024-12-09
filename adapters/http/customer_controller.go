package http

import (
	"fmt"
	"net/http"
	"os"
	"payment-layer-card-api/bootstrap"
	customers_usecase "payment-layer-card-api/usecases/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	customer_workflow "payment-layer-card-api/workflows/customer"
	"runtime"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/layer-tools/tenant"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomerController struct {
	app *bootstrap.App
}

func NewCustomerController(app *bootstrap.App) *CustomerController {
	return &CustomerController{app: app}
}

// TODO: MUDAR PARA SNAKE CASE
func ChooseCreateCustomerField(createCustomerInput *customer_dto.CreateCustomerDTOInput) *customer_dto.CreateCustomerDTOInput {
	if createCustomerInput.Address.PostalCode != "" {
		createCustomerInput.Address.PostalCodeAlt = ""
	} else {
		if createCustomerInput.Address.PostalCodeAlt != "" {
			createCustomerInput.Address.PostalCode = createCustomerInput.Address.PostalCodeAlt
			createCustomerInput.Address.PostalCodeAlt = ""
		}
	}

	return createCustomerInput
}

// @Summary Cria um cliente
// @Description Cria um cliente e o tokeniza nos gateways de pagamento
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param customer body customer_dto.CreateCustomerDTOInput true "Informações sobre o cliente"
// @Success 201 {object} customer_dto.CreateCustomerDTOOutput
// @Failure 400 {object} ErrorWrapperInternalDTO
// @Router /customers [post]
// @Security x-api-key
func (cc *CustomerController) CreateCustomer(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			runtime.Stack(buf, false)
			fmt.Printf("Pânico recuperado: %v\n", r)
			fmt.Printf("Stack trace: %s\n", string(buf))
			errorInstance := errors.NewError(errors.InternalServerError, nil, "Error interno do servidor")
			errors.EchoErrorResponse(c, *errorInstance, errorInstance.Message)
		}
	}()
	customerDTO := &customer_dto.CreateCustomerDTOInput{}
	err := c.Bind(customerDTO)
	ctx := c.Request().Context()

	if err != nil {
		errorController := errors.NewError(errors.ReadBody, err, "Erro desserializando DTO")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	//TODO: MUDAR PARA SNAKE CASE
	customerDTO = ChooseCreateCustomerField(customerDTO)

	err = cc.app.Validator.Struct(customerDTO)
	if err != nil {
		return errors.ValidationErrorResponse(c, err)
	}

	urlIntegration := os.Getenv("INTEGRATION_API_URL")
	key := os.Getenv("INTEGRATION_API_KEY")

	isAllowed, err := tenant.VerifyTenant(ctx, key, urlIntegration, []string{customerDTO.ServiceID})
	if err != nil {
		errorController := errors.NewError(errors.InternalServerError, err, "Erro ao verificar tenant")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	if !isAllowed {
		errorController := errors.NewError(errors.Unauthorized, nil, "ServiceID não permitido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, "Você nao tem permissão para criar um cliente vinculado a esse ServiceID")
	}

	createCustomerWorkflow := customer_workflow.NewCreateCustomerWorkflow(cc.app)
	createCustomerDtoOutput, errorCreateCustomerWorkflow := createCustomerWorkflow.Execute(ctx, customerDTO)

	if errorCreateCustomerWorkflow != nil {
		return errors.EchoErrorResponse(c, *errorCreateCustomerWorkflow, errorCreateCustomerWorkflow.LogMessage)
	}

	log.Info(ctx).Msg("Cliente criado com sucesso!")
	return c.JSON(http.StatusCreated, createCustomerDtoOutput)
}

// @Summary Obter um cliente pelo id
// @Description Buscar um cliente pelo identificador
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param id path string true "Identificador do cliente"
// @Success 200 {object} customer_dto.CustomerWithTokensOutputDTO
// @Failure 400 {object} ErrorWrapperInternalDTO
// @Router /customers/{id} [get]
// @Security x-api-key
func (cc *CustomerController) GetCustomerByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	uuidErr := uuid.Validate(id)
	if id == "" || uuidErr != nil {
		errorController := errors.NewError(errors.ParameterIsRequired, uuidErr, "ID campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	getCustomerByIdWorkflow := customer_workflow.NewGetCustomerByIdWorkflow(cc.app)
	customerWithTokensOutputDTO, err := getCustomerByIdWorkflow.Execute(ctx, id)

	if err != nil {
		return errors.EchoErrorResponse(c, *err, err.Message)
	}

	log.Info(ctx).Msg("Cliente obtido com sucesso pelo ID!")
	return c.JSON(http.StatusOK, customerWithTokensOutputDTO)
}

// @Summary Obter um cliente pelo documento
// @Description Buscar um cliente pelo documento
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param id path string true "Documento do cliente"
// @Success 200 {object} customer_dto.CustomerOutputDTO
// @Failure 400 {object} ErrorWrapperInternalDTO
// @Router /customers/document/{document}/serviceID/{serviceID} [get]
// @Security x-api-key
func (cc *CustomerController) GetCustomerByDocument(c echo.Context) error {
	document := c.Param("document")
	serviceID := c.Param("serviceID")
	ctx := c.Request().Context()

	if document == "" {
		errorController := errors.NewError(errors.ParameterIsRequired, nil, "CPF campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	if serviceID == "" {
		errorController := errors.NewError(errors.ParameterIsRequired, nil, "ServiceID campo vazío ou inválido")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	getCustomerByDocumentUsecase := customers_usecase.NewGetCustomerByDocument(cc.app.CustomerRepo)
	customerOutputDTO, err := getCustomerByDocumentUsecase.Execute(ctx, document, serviceID)

	if err != nil {
		if err.HttpStatus == http.StatusNotFound {
			return errors.EchoErrorResponse(c, *err, "Cliente não encontrado")
		}
		return errors.EchoErrorResponse(c, *err, "Error ao buscar cliente por documento")
	}

	log.Info(ctx).Msg("Cliente obtido com sucesso pelo Documento!")
	return c.JSON(http.StatusOK, customerOutputDTO)
}

// // @Summary Atualizar um cliente
// // @Description Atualiza um cliente
// // @Tags Customers
// // @Accept  json
// // @Produce  json
// // @Param id path string true "Identificador do cliente"
// // @Param customer body customer_dto.UpdateCustomerDTO true "Informações sobre o cliente"
// // @Success 200 {object} customer_dto.CustomerOutputDTO
// // @Failure 400 {object} ErrorWrapperInternalDTO
// // @Router /customers/{id} [put]
// // @Security x-api-key
// func (cc *CustomerController) UpdateCustomer(c echo.Context) error {
// 	ctx := c.Request().Context()
// 	id := c.Param("id")
// 	uuidErr := uuid.Validate(id)
// 	if id == "" || uuidErr != nil {
// 		errorController := errors.NewError(errors.ParameterIsRequired, "ID campo vazío ou inválido")
// 		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
// 		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
// 	}

// 	customerDTO := &customer_dto.UpdateCustomerDTO{}
// 	err := c.Bind(customerDTO)
// 	if err != nil {
// 		errorController := errors.NewError(errors.ReadBody, "Error deserializing body: "+err.Error())
// 		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
// 		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
// 	}

// 	err = cc.app.Validator.Struct(customerDTO)
// 	if err != nil {
// 		return errors.ValidationErrorResponse(c, err)
// 	}

// 	updateCustomerWorkflow := customer_workflow.NewUpdateCustomerWorkFlow(cc.app)
// 	customerOutputDTO, errResp := updateCustomerWorkflow.Execute(ctx, customerDTO, id)

// 	if errResp != nil {
// 		return errors.EchoErrorResponse(c, *errResp, errResp.Message)
// 	}

// 	log.Info(ctx).Msg("Cliente atualizado com sucesso!")
// 	return c.JSON(http.StatusOK, customerOutputDTO)
// }
