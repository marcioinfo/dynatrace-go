package card_token

import (
	"context"
	cardtoken "payment-layer-card-api/entities/card_token"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCardTokenByCardID struct {
	CardTokenRepo cardtoken.CardTokenRepositoryInterface
}

func NewGetCardTokenByCardID(cardTokenRepo cardtoken.CardTokenRepositoryInterface) *GetCardTokenByCardID {
	return &GetCardTokenByCardID{
		CardTokenRepo: cardTokenRepo,
	}
}

func (c *GetCardTokenByCardID) Execute(ctx context.Context, cardID string) ([]*cardtoken.CardToken, *errors.ErrorOutput) {
	if cardID == "" {
		return nil, errors.NewError(errors.ParameterIsRequired, nil, "ID parâmetro é obrigatório.")
	}

	cardToken, err := c.getCardTokenByCardIDInDatabase(ctx, cardID)

	if err != nil {
		return nil, err
	}

	return cardToken, nil
}

func (c *GetCardTokenByCardID) getCardTokenByCardIDInDatabase(ctx context.Context, cardID string) ([]*cardtoken.CardToken, *errors.ErrorOutput) {
	cardToken, err := c.CardTokenRepo.GetByCardID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	return cardToken, nil
}
