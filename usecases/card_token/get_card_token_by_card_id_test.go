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

func TestGetCardTokenByCardID_Execute_TokenExists(t *testing.T) {
	cardID := "sampleCardID"

	mockToken := []*card_token_entity.CardToken{{CardID: cardID, CardToken: "sampleToken"}}
	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	mockRepo.On("GetByCardID", cardID).Return(mockToken, nil)

	getToken := card_token_usecase.NewGetCardTokenByCardID(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, cardID, token[0].CardID)

}

func TestGetCardTokenByCardID_Execute_TokenNotFound(t *testing.T) {
	cardID := "sampleCardId"

	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	mockRepo.On("GetByCardID", cardID).Return(nil, sql.ErrNoRows)

	getToken := card_token_usecase.NewGetCardTokenByCardID(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID)

	assert.NotNil(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Type(), layerErrors.NotFoundError)

}

func TestGetCardTokenByCardID_Execute_RepositoryError(t *testing.T) {
	cardID := "sampleCardId"

	mockRepoErr := errors.New("database error")
	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	mockRepo.On("GetByCardID", cardID).Return(nil, mockRepoErr)

	getToken := card_token_usecase.NewGetCardTokenByCardID(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID)

	assert.NotNil(t, err)
	assert.Nil(t, token)
}

func TestGetCardTokenByCardID_Execute_InvalidID(t *testing.T) {
	cardID := ""

	mockRepo := new(card_token_entity.CardTokenRepositoryMock)
	getToken := card_token_usecase.NewGetCardTokenByCardID(mockRepo)

	ctx := context.Background()
	token, err := getToken.Execute(ctx, cardID)

	assert.NotNil(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Type(), layerErrors.ParameterIsRequired)
}
