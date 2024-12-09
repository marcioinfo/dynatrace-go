package pagarme

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"payment-layer-card-api/adapters/payment_gateways/pagarme/dtos"
	customer_pagarme_dto "payment-layer-card-api/adapters/payment_gateways/pagarme/dtos/customer"
	"payment-layer-card-api/adapters/payment_gateways/pagarme/utils"
	"payment-layer-card-api/common/types"
	customertoken_dto "payment-layer-card-api/usecases/customer_token/dtos"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"strconv"
	"time"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"go.elastic.co/apm/v2"
)

func (p *PagarmeAdapter) CreateCustomer(ctx context.Context, customer *customer_dto.CreateCustomerDTOInput, id string) (CustomerTokenOutputDTO *customertoken_dto.CustomerTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	inputToPagarme, errResp := customer_pagarme_dto.ConvertToCustomerDTO(customer, id)
	if errResp != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errResp, "Erro a converter customer para pagarme customer")
	}

	customerRetryTimeoutInMillisecond, _ := strconv.ParseInt(os.Getenv("CUSTOMER_RETRY_TIMEOUT"), 10, 64)

	contextTimeout := (time.Duration(customerRetryTimeoutInMillisecond) * time.Millisecond)
	newCtx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	req, errInstance := utils.RequestHttpWithContext(newCtx, inputToPagarme, "/customers", http.MethodPost)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer, errorMessage: "+errInstance.Error())
	}

	resp, errInstance := utils.TracingClient.Do(req)
	if errInstance != nil {
		apm.CaptureError(newCtx, errInstance).Send()
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	defer resp.Body.Close()

	body, errInstance := io.ReadAll(resp.Body)
	if errInstance != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	if resp.StatusCode != http.StatusOK {
		responseError := &dtos.ResponseErrorDTO{}
		errInstance = json.Unmarshal(body, responseError)
		if errInstance != nil {
			return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
		}
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, nil, "Erro ao criar customer, errorMessage: "+responseError.Message)
	}

	customerResponse := &customer_pagarme_dto.AddCustomerOutput{}
	errInstance = json.Unmarshal(body, customerResponse)

	if errInstance != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	return &customertoken_dto.CustomerTokenOutputDTO{
		Token:   customerResponse.CustomerToken,
		Gateway: types.PAGARME.ToString(),
	}, resp.StatusCode, nil
}
