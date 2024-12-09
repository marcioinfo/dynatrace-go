package cards

import (
	"context"
	"payment-layer-card-api/entities/cards"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type DeleteCardByID struct {
	CardRepo cards.CardRepositoryInterface
}

func NewDeleteCardByID(cardRepo cards.CardRepositoryInterface) *DeleteCardByID {
	return &DeleteCardByID{
		CardRepo: cardRepo,
	}
}

func (c *DeleteCardByID) Execute(ctx context.Context, id string) (string, *errors.ErrorOutput) {
	if id == "" {
		return "", errors.NewError(errors.DocumentIsRequired, nil, "ID é obrigatório.")
	}

	id, err := c.deleteCardByIDInDatabase(ctx, id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *DeleteCardByID) deleteCardByIDInDatabase(ctx context.Context, id string) (string, *errors.ErrorOutput) {
	id, err := c.CardRepo.DeleteByID(ctx, id)
	if err != nil {
		return "", err
	}

	return id, nil
}
