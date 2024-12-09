package cards_test

import (
	"context"
	"payment-layer-card-api/entities/cards"
	"payment-layer-card-api/mocks"
	cards_usecase "payment-layer-card-api/usecases/cards"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GetCardByIdAndCustomerIdUnitTest struct {
	cardRepository *mocks.CardRepositoryInterface
}

func NewGetCardUnitTest(t *testing.T) *GetCardByIdAndCustomerIdUnitTest {
	cardRepo := mocks.NewCardRepositoryInterface(t)

	return &GetCardByIdAndCustomerIdUnitTest{
		cardRepository: cardRepo,
	}
}

func (g *GetCardByIdAndCustomerIdUnitTest) initGetByIdAndCustomerIdMockFuncs() {
	g.cardRepository.On("GetByIDAndCustomerID", mock.Anything, mock.Anything, "f5bd24c2-cac0-4c4f-ac8f-3c5f27819a52").Return(&cards.Card{
		ID: "f5bd24c2-cac0-4c4f-ac8f-3c5f27819a52",
	}, nil)
}

func Test_GetByIdAndCustomerId_Execute(t *testing.T) {
	testTable := []struct {
		Id             string
		CustomerId     string
		ExpectedError  bool
		ExpectedOutput *cards.Card
	}{
		{
			Id:            "f5bd24c2-cac0-4c4f-ac8f-3c5f27819a52",
			CustomerId:    "f5bd24c2-cac0-4c4f-ac8f-3c5f27819a52",
			ExpectedError: false,
			ExpectedOutput: &cards.Card{
				ID: "f5bd24c2-cac0-4c4f-ac8f-3c5f27819a52",
			},
		},
		{
			Id:             "f5bd24c2-cac0-4c4f-ac8f-3c5f27819a52",
			CustomerId:     "",
			ExpectedError:  true,
			ExpectedOutput: nil,
		},
	}

	getByIdAndCustomerIdUseCase := NewGetCardUnitTest(t)
	getByIdAndCustomerIdUseCase.initGetByIdAndCustomerIdMockFuncs()
	getCard := cards_usecase.NewGetCardByIDAndCustomerID(getByIdAndCustomerIdUseCase.cardRepository)

	ctx := context.Background()
	for _, tt := range testTable {
		_, err := getCard.Execute(ctx, tt.Id, tt.CustomerId)
		boolErr := err != nil
		assert.Equal(t, tt.ExpectedError, boolErr)
	}
}
