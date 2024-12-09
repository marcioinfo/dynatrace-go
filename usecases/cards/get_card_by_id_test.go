package cards_test

import (
	"context"
	"database/sql"
	"errors"
	cards_entity "payment-layer-card-api/entities/cards"
	cards_usecase "payment-layer-card-api/usecases/cards"
	"testing"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetCardByID_Execute_CardExists(t *testing.T) {
	wantID := "sampleCardId"
	wantErr := false

	mockCard := &cards_entity.Card{ID: wantID}
	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("GetByID", wantID).Return(mockCard, nil)

	getCardByID := cards_usecase.NewGetCardByID(mockRepo)

	ctx := context.Background()
	card, err := getCardByID.Execute(ctx, wantID)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, card)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, card)
		assert.Equal(t, wantID, card.ID)
	}
}

func TestGetCardByID_Execute_CardNotFound(t *testing.T) {
	wantID := "sampleCardId"
	wantErr := true

	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("GetByID", wantID).Return(nil, sql.ErrNoRows)

	getCardByID := cards_usecase.NewGetCardByID(mockRepo)

	ctx := context.Background()
	card, err := getCardByID.Execute(ctx, wantID)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, card)
		assert.Contains(t, err.Type(), layerErrors.NotFoundError)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, card)
	}
}

func TestGetCardByID_Execute_RepositoryError(t *testing.T) {
	wantID := "sampleCardId"
	wantErr := true

	mockRepoErr := errors.New("database error")
	mockRepo := new(cards_entity.CardRepositoryMock)
	mockRepo.On("GetByID", wantID).Return(nil, mockRepoErr)

	getCardByID := cards_usecase.NewGetCardByID(mockRepo)

	ctx := context.Background()
	card, err := getCardByID.Execute(ctx, wantID)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, card)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, card)
	}
}

func TestGetCardByID_Execute_InvalidID(t *testing.T) {
	wantErr := true

	mockRepo := new(cards_entity.CardRepositoryMock)

	getCardByID := cards_usecase.NewGetCardByID(mockRepo)

	ctx := context.Background()
	card, err := getCardByID.Execute(ctx, "")

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, card)
		assert.Contains(t, err.Type(), layerErrors.ParameterIsRequired)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, card)
	}
}
