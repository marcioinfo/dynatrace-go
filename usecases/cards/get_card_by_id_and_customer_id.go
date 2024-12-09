package cards

import (
	"context"
	"payment-layer-card-api/common/helpers"
	"payment-layer-card-api/entities/cards"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCardByIDAndCustomerID struct {
	CardRepo cards.CardRepositoryInterface
}

func NewGetCardByIDAndCustomerID(cardRepo cards.CardRepositoryInterface) *GetCardByIDAndCustomerID {
	return &GetCardByIDAndCustomerID{
		CardRepo: cardRepo,
	}
}

func (c *GetCardByIDAndCustomerID) Execute(ctx context.Context, id string, customerId string) (*cards.Card, *errors.ErrorOutput) {
	if id == "" {
		return nil, errors.NewError(errors.ParameterIsRequired, nil, "ID parâmetro é obrigatório.")
	}

	if !helpers.IsValidUUID(id) || !helpers.IsValidUUID(customerId) {
		return nil, errors.NewError(errors.UUIDIsInvalid, nil, "CustomerID ou CardID está inválido.")
	}

	card, err := c.getCardByIDInDatabase(ctx, id, customerId)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (c *GetCardByIDAndCustomerID) getCardByIDInDatabase(ctx context.Context, id string, customerId string) (*cards.Card, *errors.ErrorOutput) {
	card, err := c.CardRepo.GetByIDAndCustomerID(ctx, id, customerId)
	if err != nil {
		return nil, err
	}

	return card, nil
}
