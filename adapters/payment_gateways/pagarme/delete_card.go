package pagarme

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"payment-layer-card-api/adapters/payment_gateways/pagarme/dtos"
	card_pagarme_dto "payment-layer-card-api/adapters/payment_gateways/pagarme/dtos/credit_card"
	"payment-layer-card-api/adapters/payment_gateways/pagarme/utils"
	"payment-layer-card-api/common/types"
	cardtoken_dto "payment-layer-card-api/usecases/card_token/dtos"
	"strconv"
	"time"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"go.elastic.co/apm/v2"
)

func (p *PagarmeAdapter) DeleteCard(ctx context.Context, customerId string, cardId string) (cardTokenOutputDTO *cardtoken_dto.CardTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	cardRetryTimeoutInMillisecond, _ := strconv.ParseInt(os.Getenv("CARD_RETRY_TIMEOUT"), 10, 64)

	contextTimeout := (time.Duration(cardRetryTimeoutInMillisecond) * time.Millisecond)
	newCtx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	pagarmeQuery := "customers/" + customerId + "/cards/" + cardId

	req, errInstance := utils.RequestHttpWithContext(newCtx, nil, pagarmeQuery, http.MethodDelete)
	if errInstance != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão")
	}

	resp, errInstance := utils.TracingClient.Do(req)
	if errInstance != nil {
		apm.CaptureError(newCtx, errInstance).Send()
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão")
	}

	defer resp.Body.Close()

	body, errInstance := io.ReadAll(resp.Body)
	if errInstance != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão")
	}

	if resp.StatusCode != http.StatusOK {

		responseError := &dtos.ResponseErrorDTO{}
		errInstance = json.Unmarshal(body, responseError)
		if errInstance != nil {
			return nil, resp.StatusCode, errors.NewError(errors.ExternalError, errInstance, "Erro ao deletar cartão")
		}

		errorMessages := dtos.GetAllErrorsByResponse(responseError)

		if errorMessages != "" {
			return nil, resp.StatusCode, errors.NewError(errors.ExternalError, nil, "Erro ao deletar cartão, errorMessage: "+errorMessages)
		}

		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, nil, "Erro ao deletar cartão, errorMessage: "+responseError.Message)
	}

	cardResponse := &card_pagarme_dto.AddAndTokemizeCardOutput{}

	unmarshallErr := json.Unmarshal(body, cardResponse)
	if unmarshallErr != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, unmarshallErr, "Erro ao deletar cartão")
	}

	return &cardtoken_dto.CardTokenOutputDTO{
		Token:   cardResponse.CardToken,
		Gateway: types.PAGARME.ToString(),
	}, resp.StatusCode, nil
}
