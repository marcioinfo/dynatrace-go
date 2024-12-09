package card_token

import (
	"context"
	card_token_entity "payment-layer-card-api/entities/card_token"
	payment_gateway_entity "payment-layer-card-api/entities/payment_gateway"
	card_token_dto "payment-layer-card-api/usecases/card_token/dtos"
	"strconv"
	"time"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	"github.com/adhfoundation/layer-tools/log"
)

type CreateCardToken struct {
	CardTokenRepository          card_token_entity.CardTokenRepositoryInterface
	PaymentGateway               payment_gateway_entity.PaymentGatewayInterface
	CardRetryCount               int
	CardRetryDelayInMilliseconds int
}

func NewCreateCardToken(
	cardTokenRepository card_token_entity.CardTokenRepositoryInterface,
	paymentGateway payment_gateway_entity.PaymentGatewayInterface,
	cardRetryCount int,
	cardRetryDelayInMilliseconds int,
) *CreateCardToken {

	return &CreateCardToken{
		CardTokenRepository:          cardTokenRepository,
		PaymentGateway:               paymentGateway,
		CardRetryCount:               cardRetryCount,
		CardRetryDelayInMilliseconds: cardRetryDelayInMilliseconds,
	}
}

func (cc *CreateCardToken) Execute(ctx context.Context, card_id string, input *card_token_dto.CardTokenInputDTO) (*card_token_dto.CardTokenOutputDTO, *layerErrors.ErrorOutput) {
	cardToken := card_token_entity.NewCardToken()
	cardToken.InitID()

	cardTokenDTOOutput, errorCreateCardTokenInGateway := cc.createCardTokenInGateway(ctx, input)

	if errorCreateCardTokenInGateway != nil {
		log.Error(ctx, errorCreateCardTokenInGateway.LogMessageToError()).Msgf("Erro ao criar token de cartão no gateway. Erro (%v)", errorCreateCardTokenInGateway)
		return nil, errorCreateCardTokenInGateway
	}

	cardToken.CardID = input.CardID
	cardToken.Gateway = cardTokenDTOOutput.Gateway
	cardToken.CardToken = cardTokenDTOOutput.Token

	errorInsertCardTokenInDatabase := cc.insertCardInDatabase(ctx, cardToken)
	if errorInsertCardTokenInDatabase != nil {
		log.Error(ctx, errorInsertCardTokenInDatabase.LogMessageToError()).Msgf("Erro ao criar token de cartão no banco de dados. Erro (%v)", errorInsertCardTokenInDatabase)
		return nil, errorInsertCardTokenInDatabase
	}

	return cardTokenDTOOutput, nil
}

func (cc *CreateCardToken) insertCardInDatabase(ctx context.Context, cardToken *card_token_entity.CardToken) *layerErrors.ErrorOutput {
	err := cc.CardTokenRepository.Insert(ctx, cardToken)
	if err != nil {
		return err
	}
	return nil
}

func (cc *CreateCardToken) createCardTokenInGateway(ctx context.Context, input *card_token_dto.CardTokenInputDTO) (*card_token_dto.CardTokenOutputDTO, *layerErrors.ErrorOutput) {
	var cardTokenOutputDTO *card_token_dto.CardTokenOutputDTO
	var httpStatusCode int
	var errorOnCreateCard *layerErrors.ErrorOutput

	cardRetryCount := cc.CardRetryCount
	for i := 1; i <= cardRetryCount; i++ {
		if i > 1 {
			cardRetryDelay := (time.Duration(cc.CardRetryDelayInMilliseconds) * time.Millisecond)

			log.Info(ctx).Msgf("Aguardando (%v) para a Tentativa (%v) da criação do Token de Cartão para o CardID (%v) do CustomerID (%v) no Gateway (%v)",
				cardRetryDelay,
				i,
				input.CardID,
				input.CustomerID,
				cc.PaymentGateway.GatewayName(),
			)

			time.Sleep(cardRetryDelay)
		}

		log.Info(ctx).Msgf("Tentativa (%v) da criação do Token de Cartão para o CardID (%v) do CustomerID (%v) no Gateway (%v)",
			i,
			input.CardID,
			input.CustomerID,
			cc.PaymentGateway.GatewayName(),
		)
		cardTokenOutputDTO, httpStatusCode, errorOnCreateCard = cc.PaymentGateway.CreateCard(ctx, input)

		if errorOnCreateCard == nil {
			log.Info(ctx).Msgf("Token de Cartão criado com sucesso para o CardID (%v) do CustomerID (%v) no Gateway (%v) na tentativa (%v)",
				input.CardID,
				input.CustomerID,
				cc.PaymentGateway.GatewayName(),
				i,
			)
			return cardTokenOutputDTO, nil
		}

		strHttpStatusCode := strconv.Itoa(httpStatusCode)
		log.Error(ctx, errorOnCreateCard.LogMessageToError()).Msgf("Ocorreu um erro na criação do Token de Cartão para o CardID (%v) do CustomerID (%v) no Gateway (%v) na tentativa (%v) e retornou httpStatus (%v) com o erro (%v) \n",
			input.CardID,
			input.CustomerID,
			cc.PaymentGateway.GatewayName(),
			i,
			strHttpStatusCode,
			errorOnCreateCard,
		)
		if strHttpStatusCode[0] == '4' {
			return nil, errorOnCreateCard
		}

	}

	log.Error(ctx, errorOnCreateCard.LogMessageToError()).Msgf("Não foi possivel criar o token para o cartão com CardID (%v) do CustomerID (%v) no Gateway (%v) depois de (%v) tentativas",
		input.CardID,
		input.CustomerID,
		cc.PaymentGateway.GatewayName(),
		cardRetryCount,
	)

	return nil, layerErrors.NewError(layerErrors.InternalServerError, nil, errorOnCreateCard.Error())
}
