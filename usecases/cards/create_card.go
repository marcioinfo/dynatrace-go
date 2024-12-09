package cards

import (
	"context"

	"payment-layer-card-api/entities/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CreateCard struct {
	CardRepository cards.CardRepositoryInterface
}

func NewCreateCard(cardRepository cards.CardRepositoryInterface) *CreateCard {
	return &CreateCard{
		CardRepository: cardRepository,
	}
}

func (c *CreateCard) Execute(ctx context.Context, input *card_dtos.CreateCardDTOInput) (*cards.Card, *errors.ErrorOutput) {

	card := cards.NewCardApplication()
	card.InitID()
	card.Holder = input.Holder
	card.Brand = input.Brand
	card.FirstDigits = input.Number[0:4]
	card.LastDigits = input.Number[len(input.Number)-4:]
	card.ExpMonth = input.ExpMonth
	card.ExpYear = input.ExpYear
	card.CustomerID = input.CustomerID

	card.Fingerprint = card.GenerateFingerprint(input.Number)

	valid := card.IsValid()
	if valid != nil {
		return nil, errors.NewError(errors.ValidationEntityError, valid)
	}

	err_resp := c.insertCardInDatabase(ctx, card)
	if err_resp != nil {
		return nil, err_resp
	}

	return card, nil

}

func (c *CreateCard) insertCardInDatabase(ctx context.Context, card *cards.Card) *errors.ErrorOutput {
	err := c.CardRepository.Insert(ctx, card)
	if err != nil {
		return err
	}
	return nil
}
