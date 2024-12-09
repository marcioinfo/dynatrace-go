package cards_test

import (
	"context"
	"errors"
	cards_entity "payment-layer-card-api/entities/cards"
	cards_usecase "payment-layer-card-api/usecases/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"
	"testing"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVerifyCardAlreadyExists_Execute_CardExists(t *testing.T) {
	input := card_dtos.CreateCardDTOInput{
		Number:     "1234567812345678",
		ExpMonth:   "12",
		ExpYear:    "2025",
		Brand:      "Visa",
		CustomerID: "sampleCustomerId",
		Holder:     "John Doe",
		CVV:        "123",
		BillingAddress: &card_dtos.BillingAddress{
			Line1:      "Rua",
			State:      "Estado",
			City:       "Cidade",
			Country:    "BR",
			Line2:      "Bairro",
			District:   "Distrito",
			PostalCode: "12345678",
		},
	}

	mockCard := &cards_entity.Card{
		ID:          "someId",
		CustomerID:  "78539010-7ffa-4167-bd4a-eae0ad6ca55f",
		Fingerprint: "46070d4bf934fb0d4b06d9e2c46e346944e322444900a435d7d9a95e6d7435f5",
	}
	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("GetCardByFingerprint", mock.AnythingOfType("string")).Return(mockCard, nil)

	verifyCard := cards_usecase.NewVerifyCardAlreadyExists(mockRepo)

	ctx := context.Background()
	card, err := verifyCard.Execute(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, card)

}

func TestVerifyCardAlreadyExists_Execute_CardNotFound(t *testing.T) {
	input := card_dtos.CreateCardDTOInput{
		Number:     "1234567812345678",
		ExpMonth:   "12",
		ExpYear:    "2025",
		Brand:      "Visa",
		CustomerID: "sampleCustomerId",
		Holder:     "John Doe",
		CVV:        "123",
		BillingAddress: &card_dtos.BillingAddress{
			Line1:      "Rua",
			State:      "Estado",
			City:       "Cidade",
			Country:    "BR",
			Line2:      "Bairro",
			District:   "Distrito",
			PostalCode: "12345678",
		},
	}

	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("GetCardByFingerprint", mock.AnythingOfType("string")).Return(nil, &layerErrors.ErrorOutput{Code: "DB-500", Message: "Erro interno do servidor", LogMessage: "", HttpStatus: 500})

	verifyCard := cards_usecase.NewVerifyCardAlreadyExists(mockRepo)

	ctx := context.Background()
	card, err := verifyCard.Execute(ctx, input)

	assert.Equal(t, err.Code, layerErrors.InternalServerError)
	assert.Nil(t, card)
}

func TestVerifyCardAlreadyExists_Execute_RepositoryError(t *testing.T) {
	input := card_dtos.CreateCardDTOInput{
		Number:     "1234567812345678",
		ExpMonth:   "12",
		ExpYear:    "2025",
		Brand:      "Visa",
		CustomerID: "sampleCustomerId",
		Holder:     "John Doe",
		CVV:        "123",
		BillingAddress: &card_dtos.BillingAddress{
			Line1:      "Rua",
			State:      "Estado",
			City:       "Cidade",
			Country:    "BR",
			Line2:      "Bairro",
			District:   "Distrito",
			PostalCode: "12345678",
		},
	}

	mockRepoErr := errors.New("database error")
	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("GetCardByFingerprint", mock.AnythingOfType("string")).Return(nil, mockRepoErr)

	verifyCard := cards_usecase.NewVerifyCardAlreadyExists(mockRepo)

	ctx := context.Background()
	card, err := verifyCard.Execute(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, card)
}
