package card_token_test

import (
	"context"
	"net/http"
	"payment-layer-card-api/mocks"
	card_token_usecase "payment-layer-card-api/usecases/card_token"
	card_token_dto "payment-layer-card-api/usecases/card_token/dtos"
	"testing"

	error_pkg "github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	"github.com/stretchr/testify/mock"
)

type CardTokenDeleteTest struct {
	cardTokenRepository     *mocks.CardTokenRepositoryInterface
	paymentGatewayInterface *mocks.PaymentGatewayInterface
}

func NewCardTokenUseCaseTest(t *testing.T) *CardTokenDeleteTest {
	cardTokenRepo := mocks.NewCardTokenRepositoryInterface(t)
	paymentInterface := mocks.NewPaymentGatewayInterface(t)

	return &CardTokenDeleteTest{
		cardTokenRepository:     cardTokenRepo,
		paymentGatewayInterface: paymentInterface,
	}
}

func (c *CardTokenDeleteTest) initMocks() {
	c.cardTokenRepository.On("DeleteByCardToken", mock.Anything, "1").Return(nil)
	c.cardTokenRepository.On("DeleteByCardToken", mock.Anything, "3").Return(error_pkg.NewError(error_pkg.DatabaseConnectionError, nil, "Error"))
	c.paymentGatewayInterface.On("GatewayName").Return("pagarme")
	c.paymentGatewayInterface.On("DeleteCard", mock.Anything, "1", "1").Return(
		&card_token_dto.CardTokenOutputDTO{
			Gateway: "pagarme",
			Token:   "1",
		},
		http.StatusOK,
		nil,
	)
	c.paymentGatewayInterface.On("DeleteCard", mock.Anything, "2", "2").Return(
		nil,
		http.StatusBadRequest,
		error_pkg.NewPaymentLayerError(error_pkg.InternalServerError, "Error"),
	)
	c.paymentGatewayInterface.On("DeleteCard", mock.Anything, "3", "3").Return(
		&card_token_dto.CardTokenOutputDTO{
			Gateway: "pagarme",
			Token:   "3",
		},
		http.StatusOK,
		nil,
	)
}

func Test_DeleteCardToken_Execute(t *testing.T) {
	ctx := context.Background()

	deleteCardTokenTest := NewCardTokenUseCaseTest(t)
	deleteCardTokenTest.initMocks()

	testTable := []struct {
		name                string
		cardIdInGateway     string
		customerIdInGateway string
		ExpectedError       bool
		ExpectedOutput      *card_token_dto.CardTokenOutputDTO
	}{
		{
			name:                "successful case",
			cardIdInGateway:     "1",
			customerIdInGateway: "1",
			ExpectedError:       false,
			ExpectedOutput: &card_token_dto.CardTokenOutputDTO{
				Gateway: "pagarme",
				Token:   "1",
			},
		},
		{
			name:                "error in gateway case",
			cardIdInGateway:     "2",
			customerIdInGateway: "2",
			ExpectedError:       true,
			ExpectedOutput:      nil,
		},
		{
			name:                "error in database case",
			cardIdInGateway:     "3",
			customerIdInGateway: "3",
			ExpectedError:       true,
			ExpectedOutput:      nil,
		},
	}

	deleteCardToken := card_token_usecase.NewDeleteCardToken(
		deleteCardTokenTest.cardTokenRepository,
		deleteCardTokenTest.paymentGatewayInterface,
	)

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			_, err := deleteCardToken.Execute(ctx, tt.customerIdInGateway, tt.cardIdInGateway)

			hasErr := err != nil

			if hasErr != tt.ExpectedError {
				t.Fail()
			}
		})

	}
}
