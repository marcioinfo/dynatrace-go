package cards

import (
	"context"
	"payment-layer-card-api/entities/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type VerifyCardAlreadyExists struct {
	CardRepo cards.CardRepositoryInterface
}

func NewVerifyCardAlreadyExists(cardRepo cards.CardRepositoryInterface) *VerifyCardAlreadyExists {
	return &VerifyCardAlreadyExists{
		CardRepo: cardRepo,
	}
}

func (c *VerifyCardAlreadyExists) Execute(ctx context.Context, input card_dtos.CreateCardDTOInput) (*cards.Card, *errors.ErrorOutput) {
	card := cards.Card{}

	card.Brand = input.Brand
	card.CustomerID = input.CustomerID
	card.FirstDigits = input.Number[0:4]
	card.LastDigits = input.Number[len(input.Number)-4:]
	card.ExpMonth = input.ExpMonth
	card.ExpYear = input.ExpYear
	card.Holder = input.Holder

	card.Fingerprint = card.GenerateFingerprint(input.Number)

	cardExists, err := c.getCardByFingerprintInDatabase(ctx, card.Fingerprint)

	if err != nil && err.Code == errors.NotFoundError {
		log.Info(ctx).Msgf("O cartão do CustomerID (%v) não foi encontrado para o Fingerprint (%v)",
			card.CustomerID,
			card.Fingerprint,
		)
		return nil, nil
	}

	if err != nil && err.Code != errors.NotFoundError {
		log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao pesquisar o cartão do CustomerID (%v) pelo Fingerprint (%v). Erro (%v)",
			card.CustomerID,
			card.Fingerprint,
			err,
		)
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	log.Info(ctx).Msgf("O cartão com CardID (%v) do CustomerID (%v) foi encontrado pelo Fingerprint (%v)",
		cardExists.ID,
		cardExists.CustomerID,
		cardExists.Fingerprint,
	)

	return cardExists, nil
}

func (c *VerifyCardAlreadyExists) getCardByFingerprintInDatabase(ctx context.Context, fingerprint string) (*cards.Card, *errors.ErrorOutput) {
	card, err := c.CardRepo.GetCardByFingerprint(ctx, fingerprint)
	if err != nil {
		return nil, err
	}
	return card, nil
}
