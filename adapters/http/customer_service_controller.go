package http

import (
	"payment-layer-card-api/bootstrap"
	customerservice_usecase "payment-layer-card-api/usecases/customer_service"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomerServiceController struct {
	app *bootstrap.App
}

func NewCustomerServiceController(app *bootstrap.App) *CustomerServiceController {
	return &CustomerServiceController{app: app}
}

func (csc *CustomerServiceController) GetCustomerServiceByServiceIdAndCustomerId(c echo.Context) error {
	ctx := c.Request().Context()

	serviceId := c.Param("serviceId")
	customerId := c.Param("customerId")

	if serviceId == "" || customerId == "" {
		errorController := errors.NewError(errors.ParameterIsRequired, nil, "Customer ID e Service ID são obrigatórios")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	errUUID := uuid.Validate(serviceId)
	if errUUID != nil {
		errorController := errors.NewControllerError(errors.BadRequest, "Service ID is not a valid UUID")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	errUUID = uuid.Validate(customerId)
	if errUUID != nil {
		errorController := errors.NewControllerError(errors.BadRequest, "Customer ID is not a valid UUID")
		log.Error(ctx, errorController.LogMessageToError()).Msg(errorController.Message)
		return errors.EchoErrorResponse(c, *errorController, errorController.Message)
	}

	getCustomerServiceByCustomerAndServiceID := customerservice_usecase.NewGetCustomerServiceByCustomerAndServiceID(csc.app.CustomerServiceRepo)
	customerService, errorOutput := getCustomerServiceByCustomerAndServiceID.Execute(ctx, customerId, serviceId)
	if errorOutput != nil {
		return errors.EchoErrorResponse(c, *errorOutput, errorOutput.LogMessage)
	}

	return c.JSON(200, customerService)
}
