package card_token

import (
	"context"
	"errors"

	card_token_entity "payment-layer-card-api/entities/card_token"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

var (
	ErrorCardIdRequired    = errors.New("card id é obrigatório")
	ErrorGatewayIdRequired = errors.New("gateway é obrigatório")
	ErrorCardTokenNotFound = errors.New("card token não encontrado")
)

type GetCardTokenByCardIDAndGateway struct {
	CardTokenRepository card_token_entity.CardTokenRepositoryInterface
}

func NewGetCardTokenByCardIDAndGateway(cardTokenRepository card_token_entity.CardTokenRepositoryInterface) *GetCardTokenByCardIDAndGateway {
	return &GetCardTokenByCardIDAndGateway{
		CardTokenRepository: cardTokenRepository,
	}
}

func (g *GetCardTokenByCardIDAndGateway) Execute(ctx context.Context, cardId string, gateway string) (*card_token_entity.CardToken, *layerErrors.ErrorOutput) {
	if cardId == "" {
		return nil, layerErrors.NewError(layerErrors.ParameterIsRequired, ErrorCardIdRequired)
	}

	if gateway == "" {
		return nil, layerErrors.NewError(layerErrors.WithoutGateway, ErrorGatewayIdRequired)
	}

	cardToken, err := g.getCardTokenByCardIDAndGateway(ctx, cardId, gateway)

	if err != nil {
		return nil, err
	}

	return cardToken, nil
}

func (g *GetCardTokenByCardIDAndGateway) getCardTokenByCardIDAndGateway(ctx context.Context, cardId string, gateway string) (*card_token_entity.CardToken, *layerErrors.ErrorOutput) {
	cardToken, err := g.CardTokenRepository.GetByCardIDAndGateway(ctx, cardId, gateway)
	if err != nil {
		return nil, err
	}
	return cardToken, nil
}
