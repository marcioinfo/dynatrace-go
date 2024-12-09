package card_token_test

import (
	"context"
	"database/sql"
	"errors"
	card_token_entity "payment-layer-card-api/entities/card_token"
	card_token_usecase "payment-layer-card-api/usecases/card_token"
	"testing"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestGetCardTokenByCardIDAndGateway_Execute_TokenExists(t *testing.T) {
	cardID := "sampleCardID"
	gateway := "pagarme"

	mockToken := &card_token_entity.CardToken{CardID: cardID, CardToken: "sampleToken"}
	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	mockRepo.On("GetByCardIDAndGateway", cardID, gateway).Return(mockToken, nil)

	getToken := card_token_usecase.NewGetCardTokenByCardIDAndGateway(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID, gateway)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, cardID, token.CardID)
}

func TestGetCardTokenByCardIDAndGateway_Execute_TokenNotFound(t *testing.T) {
	cardID := "sampleCardId"
	gateway := "pagarme"

	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	mockRepo.On("GetByCardIDAndGateway", cardID, gateway).Return(nil, sql.ErrNoRows)

	getToken := card_token_usecase.NewGetCardTokenByCardIDAndGateway(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID, gateway)

	assert.NotNil(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Code, layerErrors.NotFoundError)

}

func TestGetCardTokenByCardIDAndGateway_Execute_RepositoryError(t *testing.T) {
	cardID := "sampleCardId"
	gateway := "pagarme"

	mockRepoErr := errors.New("database error")
	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	mockRepo.On("GetByCardIDAndGateway", cardID, gateway).Return(nil, mockRepoErr)

	getToken := card_token_usecase.NewGetCardTokenByCardIDAndGateway(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID, gateway)

	assert.NotNil(t, err)
	assert.Nil(t, token)
}

func TestGetCardTokenByCardIDAndGateway_Execute_Invalid_CardID(t *testing.T) {
	cardID := ""
	gateway := "pagarme"

	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	getToken := card_token_usecase.NewGetCardTokenByCardIDAndGateway(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID, gateway)

	assert.NotNil(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Type(), layerErrors.ParameterIsRequired)
}

func TestGetCardTokenByCardIDAndGateway_Execute_Invalid_Gateway(t *testing.T) {
	cardID := "sampleCardId"
	gateway := ""

	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	getToken := card_token_usecase.NewGetCardTokenByCardIDAndGateway(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID, gateway)

	assert.NotNil(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Type(), layerErrors.WithoutGateway)
}
