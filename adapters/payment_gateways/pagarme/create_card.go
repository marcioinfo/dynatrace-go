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

func (p *PagarmeAdapter) CreateCard(ctx context.Context, input *cardtoken_dto.CardTokenInputDTO) (*cardtoken_dto.CardTokenOutputDTO, int, *errors.ErrorOutput) {
	inputToPagarme := card_pagarme_dto.ConvertToCreditCardDTO(input)

	cardRetryTimeoutInMillisecond, _ := strconv.ParseInt(os.Getenv("CARD_RETRY_TIMEOUT"), 10, 64)

	contextTimeout := (time.Duration(cardRetryTimeoutInMillisecond) * time.Millisecond)
	newCtx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	pagarmeQuery := "customers/" + input.CustomerToken + "/cards"

	req, err := utils.RequestHttpWithContext(newCtx, inputToPagarme, pagarmeQuery, http.MethodPost)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, err, "Erro ao criar cartão")
	}

	resp, err := utils.TracingClient.Do(req)
	if err != nil {
		apm.CaptureError(newCtx, err).Send()
		return nil, http.StatusInternalServerError, errors.NewError(errors.ExternalError, err, "Erro ao criar cartão")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, err, "Erro ao criar cartão")
	}

	errPagarme := utils.CheckForErrors(body, resp.StatusCode)
	if errPagarme != nil {
		return nil, resp.StatusCode, errPagarme
	}

	if resp.StatusCode != http.StatusOK {
		responseError := &dtos.ResponseErrorDTO{}
		err = json.Unmarshal(body, responseError)
		if err != nil {
			return nil, resp.StatusCode, errors.NewError(errors.ExternalError, err, "Erro ao criar cartão")
		}

		errorMessages := dtos.GetAllErrorsByResponse(responseError)
		if errorMessages != "" {
			return nil, resp.StatusCode, errors.NewError(errors.ExternalError, err, "Erro ao criar cartão, errorMessage: "+errorMessages)

		}

		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, err, "Erro ao criar cartão, errorMessage: "+responseError.Message)
	}

	cardResponse := &card_pagarme_dto.AddAndTokemizeCardOutput{}
	err = json.Unmarshal(body, cardResponse)
	if err != nil {
		return nil, resp.StatusCode, errors.NewError(errors.ExternalError, err, "Erro ao criar cartão")
	}

	return &cardtoken_dto.CardTokenOutputDTO{
		Token:   cardResponse.CardToken,
		Gateway: types.PAGARME.ToString(),
	}, resp.StatusCode, nil
}
