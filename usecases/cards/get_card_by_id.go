package cards

import (
	"context"
	"payment-layer-card-api/entities/cards"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCardByID struct {
	CardRepo cards.CardRepositoryInterface
}

func NewGetCardByID(cardRepo cards.CardRepositoryInterface) *GetCardByID {
	return &GetCardByID{
		CardRepo: cardRepo,
	}
}

func (c *GetCardByID) Execute(ctx context.Context, id string) (*cards.Card, *errors.ErrorOutput) {
	if id == "" {
		return nil, errors.NewError(errors.ParameterIsRequired, nil, "ID parâmetro é obrigatório.")
	}

	card, err := c.getCardByIDInDatabase(ctx, id)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (c *GetCardByID) getCardByIDInDatabase(ctx context.Context, id string) (*cards.Card, *errors.ErrorOutput) {
	card, err := c.CardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return card, nil
}
