package rede

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	card_rede_dto "payment-layer-card-api/adapters/payment_gateways/rede/credit_card"
	"payment-layer-card-api/adapters/payment_gateways/rede/utils"
	"payment-layer-card-api/common/types"
	cardtoken_dto "payment-layer-card-api/usecases/card_token/dtos"
	"strconv"
	"time"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

func (r *RedeAdapter) CreateCard(ctx context.Context, input *cardtoken_dto.CardTokenInputDTO) (cardTokenOutputDTO *cardtoken_dto.CardTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	createCardRequest, errInstance := card_rede_dto.ConvertToZeroDollarCard(input)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar cartão na Rede")
	}

	requestXml, errInstance := xml.Marshal(createCardRequest)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar cartão na Rede")
	}

	xmlBody := utils.BuildTransactionXmlBody(string(requestXml))

	cardRetryTimeoutInMillisecond, _ := strconv.ParseInt(os.Getenv("CARD_RETRY_TIMEOUT"), 10, 64)

	contextTimeout := (time.Duration(cardRetryTimeoutInMillisecond) * time.Millisecond)
	newCtx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	req, errInstance := utils.RequestHttpWithContext(newCtx, xmlBody, "/postXML")
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar cartão na Rede")
	}

	client := &http.Client{}
	resp, errInstance := client.Do(req)

	if errInstance != nil && resp == nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar cartão na Rede")
	}

	if errInstance != nil && resp != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar cartão na Rede")
	}

	defer resp.Body.Close()

	body, errInstance := io.ReadAll(resp.Body)
	if errInstance != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar cartão na Rede")
	}

	createResponse := &card_rede_dto.AddAndTokemizeCardWithZeroDollar{}

	errInstance = xml.Unmarshal(body, createResponse)
	if errInstance != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao criar cartão na Rede")
	}

	if createResponse.ErrorCode != "0" && createResponse.ErrorMessage != "" {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, nil, "Erro ao criar cartão na Rede: "+createResponse.ErrorMessage)
	}

	if createResponse.ResponseMessage == "DECLINED" {
		return nil, http.StatusPreconditionFailed, errors.NewError(errors.CartaoRestricao, nil, "Erro ao criar cartao na Rede: Cartao nao aceito")
	}

	return &cardtoken_dto.CardTokenOutputDTO{
		Token:   createResponse.SaveFile.Token,
		Gateway: types.REDE.ToString(),
	}, resp.StatusCode, nil
}
