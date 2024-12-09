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
	"payment-layer-card-api/entities/customers"
	customer_token_dto "payment-layer-card-api/usecases/customer_token/dtos"
	"strconv"
	"time"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

func (r *RedeAdapter) UpdateCustomer(ctx context.Context, customer *customers.Customer, id string) (*customer_token_dto.CustomerTokenOutputDTO, int, *errors.ErrorOutput) {
	inputToRede := customer_rede_dto.ConvertToUpdateCustomerDTO(customer, id)

	customerRetryTimeoutInMillisecond, _ := strconv.ParseInt(os.Getenv("CUSTOMER_RETRY_TIMEOUT"), 10, 64)

	contextTimeout := (time.Duration(customerRetryTimeoutInMillisecond) * time.Millisecond)
	newCtx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	requestXml, errXMl := xml.Marshal(inputToRede)
	if errXMl != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errXMl, "Erro ao atualizar customer")
	}

	xmlBody := utils.BuildApiRequestXmlBody(string(requestXml), "update-consumer")
	req, err := utils.RequestHttpWithContext(newCtx, xmlBody, "/postAPI")
	if err != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, err, "Erro ao atualizar customer")
	}

	resp, err := utils.TracingClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, err, "Erro ao atualizar customer")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, err, "Erro ao atualizar customer")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, nil, "Erro ao atualizar customer, errorMessage: "+string(body))
	}

	updateResponde := &customer_rede_dto.ConsumerResponseOutput{}

	err = xml.Unmarshal(body, updateResponde)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, err, "Erro ao atualizar customer")
	}

	if updateResponde.ErrorCode != "0" {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, nil, "Erro ao atualizar customer, errorMessage: "+updateResponde.ErrorMessage)
	}

	return &customer_token_dto.CustomerTokenOutputDTO{
		Gateway: types.REDE.ToString(),
		Token:   id,
	}, resp.StatusCode, nil

}
