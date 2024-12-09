package cards_test

import (
	"context"
	"errors"
	cards_entity "payment-layer-card-api/entities/cards"
	cards_usecase "payment-layer-card-api/usecases/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCard_Execute_ValidInput(t *testing.T) {
	input := &card_dtos.CreateCardDTOInput{
		CustomerID: "sampleCustomerId",
		Holder:     "John Doe",
		Brand:      "Visa",
		Number:     "1234123412341234",
		CVV:        "123",
		ExpMonth:   "09",
		ExpYear:    "2025",
		BillingAddress: &card_dtos.BillingAddress{
			Line1:      "Rua",
			Line2:      "Bairro",
			District:   "Distrito",
			City:       "Cidade",
			State:      "Estado",
			Country:    "BR",
			PostalCode: "12345678",
		},
	}
	mockRepoErr := error(nil)
	wantErr := false

	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("Insert", mock.AnythingOfType("*cards.Card")).Return(mockRepoErr)

	createCard := cards_usecase.NewCreateCard(mockRepo)

	ctx := context.Background()
	output, err := createCard.Execute(ctx, input)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, output)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.ID)
	}
}

func TestCreateCard_Execute_ShouldWorkWithEmptyBrand(t *testing.T) {
	input := &card_dtos.CreateCardDTOInput{
		CustomerID: "sampleCustomerId",
		Holder:     "John Doe",
		Brand:      "",
		Number:     "1234123412341234",
		CVV:        "123",
		ExpMonth:   "09",
		ExpYear:    "2025",
		BillingAddress: &card_dtos.BillingAddress{
			Line1:      "Rua",
			Line2:      "Bairro",
			District:   "Distrito",
			City:       "Cidade",
			State:      "Estado",
			Country:    "BR",
			PostalCode: "12345678",
		},
	}
	mockRepoErr := error(nil)
	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("Insert", mock.AnythingOfType("*cards.Card")).Return(mockRepoErr)

	createCard := cards_usecase.NewCreateCard(mockRepo)

	ctx := context.Background()
	output, err := createCard.Execute(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
}

func TestCreateCard_Execute_RepositoryError(t *testing.T) {
	input := &card_dtos.CreateCardDTOInput{
		CustomerID: "sampleCustomerId",
		Holder:     "John Doe",
		Brand:      "Visa",
		Number:     "1234123412341234",
		CVV:        "123",
		ExpMonth:   "09",
		ExpYear:    "2025",
		BillingAddress: &card_dtos.BillingAddress{
			Line1:      "Rua",
			Line2:      "Bairro",
			District:   "Distrito",
			City:       "Cidade",
			State:      "Estado",
			Country:    "BR",
			PostalCode: "12345678",
		},
	}
	mockRepoErr := errors.New("database error")
	wantErr := true

	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("Insert", mock.AnythingOfType("*cards.Card")).Return(mockRepoErr)

	createCard := cards_usecase.NewCreateCard(mockRepo)

	ctx := context.Background()
	output, err := createCard.Execute(ctx, input)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, output)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.ID)
	}
}
