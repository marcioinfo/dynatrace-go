package pagarme

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"payment-layer-card-api/adapters/payment_gateways/pagarme/dtos"
	customer_pagarme_dto "payment-layer-card-api/adapters/payment_gateways/pagarme/dtos/customer"
	"payment-layer-card-api/adapters/payment_gateways/pagarme/utils"
	"payment-layer-card-api/common/types"
	"payment-layer-card-api/entities/customers"
	customertoken_dto "payment-layer-card-api/usecases/customer_token/dtos"
	"strconv"
	"time"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"go.elastic.co/apm/v2"
)

func (p *PagarmeAdapter) UpdateCustomer(ctx context.Context, customer *customers.Customer, id string) (CustomerTokenOutputDTO *customertoken_dto.CustomerTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	inputToPagarme := customer_pagarme_dto.ConvertToEditCustomerDTO(customer)

	customerRetryTimeoutInMillisecond, _ := strconv.ParseInt(os.Getenv("CUSTOMER_RETRY_TIMEOUT"), 10, 64)

	contextTimeout := (time.Duration(customerRetryTimeoutInMillisecond) * time.Millisecond)
	newCtx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	req, errResp := utils.RequestHttpWithContext(newCtx, inputToPagarme, "/customers/"+id, http.MethodPut)
	if errResp != nil {
		fmt.Println("Opa caralho")
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errResp, "Erro ao fazer update do customer")
	}

	resp, errResp := utils.TracingClient.Do(req)
	if errResp != nil {
		apm.CaptureError(newCtx, errResp).Send()
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errResp, "Erro ao fazer update do customer")
	}

	defer resp.Body.Close()

	body, errResp := io.ReadAll(resp.Body)
	if errResp != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errResp, "Erro ao fazer update do customer")
	}

	if resp.StatusCode != http.StatusOK {
		responseError := &dtos.ResponseErrorDTO{}
		errResp = json.Unmarshal(body, responseError)
		if errResp != nil {
			return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errResp, "Erro ao fazer update do customer")
		}
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errResp, "Erro ao fazer update do customer, errorMessage: "+responseError.Message)
	}

	customerResponse := &customer_pagarme_dto.AddCustomerOutput{}
	errResp = json.Unmarshal(body, customerResponse)

	if errResp != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errResp, "Erro ao fazer update do customer")
	}

	return &customertoken_dto.CustomerTokenOutputDTO{
		Token:   customerResponse.CustomerToken,
		Gateway: types.PAGARME.ToString(),
	}, resp.StatusCode, nil
}
