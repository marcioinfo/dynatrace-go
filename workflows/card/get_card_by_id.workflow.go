package card

import (
	"context"
	"payment-layer-card-api/bootstrap"

	cardtoken_usecase "payment-layer-card-api/usecases/card_token"
	cards_usecase "payment-layer-card-api/usecases/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetCardByIDWorkflow struct {
	app *bootstrap.App
}

func NewGetCardByIdWorkflow(app *bootstrap.App) *GetCardByIDWorkflow {
	return &GetCardByIDWorkflow{app: app}
}

func (gcw *GetCardByIDWorkflow) Execute(ctx context.Context, id string) (*card_dtos.GetCardWithTokensDTO, *errors.ErrorOutput) {
	getCardByIdUsecase := cards_usecase.NewGetCardByID(gcw.app.CardRepo)

	card, err := getCardByIdUsecase.Execute(ctx, id)

	if err != nil {
		return nil, err
	}

	cardTokens, err := gcw.getCardTokens(ctx, card.ID)
	if err != nil {
		return nil, err
	}

	return &card_dtos.GetCardWithTokensDTO{
		ID:          card.ID,
		CustomerID:  card.CustomerID,
		Holder:      card.Holder,
		Brand:       card.Brand,
		FirstDigits: card.FirstDigits,
		LastDigits:  card.LastDigits,
		ExpMonth:    card.ExpMonth,
		ExpYear:     card.ExpYear,
		Tokens:      cardTokens,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}, nil
}

func (gcw *GetCardByIDWorkflow) getCardTokens(ctx context.Context, cardID string) ([]*card_dtos.CardToken, *errors.ErrorOutput) {
	getCardTokenUsecase := cardtoken_usecase.NewGetCardTokenByCardID(gcw.app.CardTokenRepo)

	cardTokens, err := getCardTokenUsecase.Execute(ctx, cardID)
	if err != nil {
		if err.Code == errors.DocumentIsRequired {
			return nil, errors.NewError(errors.CardNotFound, nil, "Não a token para esse cartão.")

		}
		return nil, err
	}

	var cardTokensOutput []*card_dtos.CardToken
	for _, cardTokenIt := range cardTokens {
		cardTokensOutput = append(cardTokensOutput, &card_dtos.CardToken{
			Token:   cardTokenIt.CardToken,
			Gateway: cardTokenIt.Gateway,
		})
	}

	return cardTokensOutput, nil
}
