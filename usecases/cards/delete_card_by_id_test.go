package cards_test

import (
	"context"
	"payment-layer-card-api/mocks"
	cards_usecase "payment-layer-card-api/usecases/cards"
	"testing"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DeleteCardUnitTest struct {
	cardRepository *mocks.CardRepositoryInterface
}

func (d *DeleteCardUnitTest) initDeleteMockFuncs() {
	d.cardRepository.On("DeleteByID", mock.Anything, "1").Return("1", nil)
	d.cardRepository.On("DeleteByID", mock.Anything, "2").Return("", errors.NewPaymentLayerError(errors.DocumentIsRequired, "ID é obrigatório."))
}

func NewDeleteCardUseCaseTest(t *testing.T) *DeleteCardUnitTest {
	cardRepo := mocks.NewCardRepositoryInterface(t)

	return &DeleteCardUnitTest{
		cardRepository: cardRepo,
	}
}

func Test_DeleteCard_Execute(t *testing.T) {
	testTable := []struct{
		Id string
		CustomerId string
		ExpectedError bool
		ExpectedOutput string
	}{
		{
			Id: "1",
			CustomerId: "1",
			ExpectedError: false,
			ExpectedOutput: "1",
		},
		{
			Id: "",
			CustomerId: "1",
			ExpectedError: true,
			ExpectedOutput: "",
		},
		{
			Id: "2",
			CustomerId: "2",
			ExpectedError: true,
			ExpectedOutput: "",
		},
	}


	cardUsecase := NewDeleteCardUseCaseTest(t)
	cardUsecase.initDeleteMockFuncs()
	deleteCard := cards_usecase.NewDeleteCardByID(cardUsecase.cardRepository)

	ctx := context.Background()
	for _, tt := range testTable {
		output, err := deleteCard.Execute(ctx, tt.Id)
		hasErr := err != nil
		assert.Equal(t, tt.ExpectedError, hasErr)
		assert.Equal(t, tt.ExpectedOutput, output)
	}
}
