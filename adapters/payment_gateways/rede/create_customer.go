package rede

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	customer_rede_dto "payment-layer-card-api/adapters/payment_gateways/rede/dtos/customer"
	"payment-layer-card-api/adapters/payment_gateways/rede/utils"
	"payment-layer-card-api/common/types"
	customertoken_dto "payment-layer-card-api/usecases/customer_token/dtos"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"strconv"
	"time"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

func (r *RedeAdapter) CreateCustomer(ctx context.Context, customer *customer_dto.CreateCustomerDTOInput, id string) (customerTokenOutputDTO *customertoken_dto.CustomerTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	createCustomerRequest := customer_rede_dto.ConvertToCreateCustomerDTO(customer, id)

	requestXml, errInstance := xml.Marshal(createCustomerRequest)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao converter customer para rede customer")
	}

	xmlBody := utils.BuildApiRequestXmlBody(string(requestXml), "add-consumer")

	customerRetryTimeoutInMillisecond, _ := strconv.ParseInt(os.Getenv("CUSTOMER_RETRY_TIMEOUT"), 10, 64)

	contextTimeout := (time.Duration(customerRetryTimeoutInMillisecond) * time.Millisecond)
	newCtx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	req, errInstance := utils.RequestHttpWithContext(newCtx, xmlBody, "/postAPI")
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	client := &http.Client{}
	resp, errInstance := client.Do(req)

	if errInstance != nil && resp == nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	if errInstance != nil && resp != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	defer resp.Body.Close()

	body, errInstance := io.ReadAll(resp.Body)
	if errInstance != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	createResponse := &customer_rede_dto.ConsumerResponseOutput{}

	errInstance = xml.Unmarshal(body, createResponse)
	if errInstance != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar customer")
	}

	if createResponse.ErrorCode != "0" {
		return &customertoken_dto.CustomerTokenOutputDTO{}, resp.StatusCode, errors.NewError(errors.ExternalError, nil, "Erro ao criar customer, errorMessage: "+createResponse.ErrorMessage)
	}

	return &customertoken_dto.CustomerTokenOutputDTO{
		Token:   createResponse.Result.CustomerID,
		Gateway: types.REDE.ToString(),
	}, resp.StatusCode, nil
}
