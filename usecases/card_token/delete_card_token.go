package card_token

import (
	"context"
	card_token_entity "payment-layer-card-api/entities/card_token"
	payment_gateway_entity "payment-layer-card-api/entities/payment_gateway"
	card_token_dto "payment-layer-card-api/usecases/card_token/dtos"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type DeleteCardToken struct {
	CardTokenRepository card_token_entity.CardTokenRepositoryInterface
	PaymentGateway      payment_gateway_entity.PaymentGatewayInterface
}

func NewDeleteCardToken(
	cardTokenRepository card_token_entity.CardTokenRepositoryInterface,
	paymentGateway payment_gateway_entity.PaymentGatewayInterface,
) *DeleteCardToken {

	return &DeleteCardToken{
		CardTokenRepository: cardTokenRepository,
		PaymentGateway:      paymentGateway,
	}

}

func (cc *DeleteCardToken) Execute(ctx context.Context, customerIdInGateway string, cardIdInGateway string) (*card_token_dto.CardTokenOutputDTO, *errors.ErrorOutput) {
	cardTokenDTOOutput, errorDeleteCardTokenInGateway := cc.deleteCardTokenInGateway(ctx, customerIdInGateway, cardIdInGateway)
	if errorDeleteCardTokenInGateway != nil {
		log.Error(ctx, errorDeleteCardTokenInGateway.LogMessageToError()).Msgf("Erro ao deletar token de cartão no gateway. Erro (%v)", errorDeleteCardTokenInGateway)
		return nil, errorDeleteCardTokenInGateway
	}

	errorInsertCardTokenInDatabase := cc.deleteCardTokenInDatabase(ctx, cardIdInGateway)
	if errorInsertCardTokenInDatabase != nil {
		log.Error(ctx, errorInsertCardTokenInDatabase.LogMessageToError()).Msgf("Erro ao deletar token de cartão no banco de dados. Erro (%v)", errorInsertCardTokenInDatabase)
		return nil, errorInsertCardTokenInDatabase
	}

	return cardTokenDTOOutput, nil
}

func (cc *DeleteCardToken) deleteCardTokenInDatabase(ctx context.Context, cardToken string) *errors.ErrorOutput {
	err := cc.CardTokenRepository.DeleteByCardToken(ctx, cardToken)
	if err != nil {
		return err
	}
	return nil
}

func (cc *DeleteCardToken) deleteCardTokenInGateway(ctx context.Context, customerId string, cardId string) (*card_token_dto.CardTokenOutputDTO, *errors.ErrorOutput) {
	var errorOnDeleteCard *errors.ErrorOutput

	log.Info(ctx).Msgf("Tentativa de deleção do Token de Cartão para o CardID (%v) do CustomerID (%v) no Gateway (%v)",
		cardId,
		customerId,
		cc.PaymentGateway.GatewayName(),
	)

	cardTokenOutputDTO, _, errorOnDeleteCard := cc.PaymentGateway.DeleteCard(ctx, customerId, cardId)

	if errorOnDeleteCard == nil {
		log.Info(ctx).Msgf("Token de Cartão deletado com sucesso para o CardID (%v) do CustomerID (%v) no Gateway (%v)",
			cardId,
			customerId,
			cc.PaymentGateway.GatewayName(),
		)
		return cardTokenOutputDTO, nil
	}

	log.Error(ctx, errorOnDeleteCard).Msgf("Não foi possivel deletar o token para o cartão com CardID (%v) do CustomerID (%v) no Gateway (%v)",
		cardId,
		customerId,
		cc.PaymentGateway.GatewayName(),
	)

	return nil, errors.NewGatewayError(errors.InvalidCard, errorOnDeleteCard.Error())
}
