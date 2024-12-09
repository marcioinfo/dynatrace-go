package rede

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	card_rede_dto "payment-layer-card-api/adapters/payment_gateways/rede/credit_card"
	customer_rede_dto "payment-layer-card-api/adapters/payment_gateways/rede/dtos/customer"
	"payment-layer-card-api/adapters/payment_gateways/rede/utils"
	"payment-layer-card-api/common/types"
	cardtoken_dto "payment-layer-card-api/usecases/card_token/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

func (r *RedeAdapter) DeleteCard(ctx context.Context, customerId string, cardId string) (cardTokenOutputDTO *cardtoken_dto.CardTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	deleteCardDto := &card_rede_dto.DeleteCardDTO{
		CustomerID: customerId,
		Token:      cardId,
	}

	requestXml, errInstance := xml.Marshal(deleteCardDto)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao converter cartão para cartão da Rede")
	}

	xmlBody := utils.BuildApiRequestXmlBody(string(requestXml), "delete-card-onfile")

	req, errInstance := utils.RequestHttpWithContext(ctx, xmlBody, "/postAPI")
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão")
	}

	client := &http.Client{}
	resp, errInstance := client.Do(req)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão, errorMessage: "+errInstance.Error())
	}

	defer resp.Body.Close()

	body, errInstance := io.ReadAll(resp.Body)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão")
	}

	deleteResponse := &customer_rede_dto.ConsumerResponseOutput{}
	errInstance = xml.Unmarshal(body, deleteResponse)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão")
	}

	if deleteResponse.ErrorCode != "0" {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, nil, "Erro ao deletar cartão, errorMessage: "+deleteResponse.ErrorMessage)
	}

	return &cardtoken_dto.CardTokenOutputDTO{
		Token:   cardId,
		Gateway: types.REDE.ToString(),
	}, resp.StatusCode, nil
}
