package card_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	adapters "payment-layer-card-api/adapters/http"
	"payment-layer-card-api/bootstrap"
	"payment-layer-card-api/entities/card_token"
	"payment-layer-card-api/entities/cards"
	customerservice "payment-layer-card-api/entities/customer_service"
	paymentgateway "payment-layer-card-api/entities/payment_gateway"
	"payment-layer-card-api/mocks"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"
	"strings"
	"testing"

	errors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/labstack/echo/v4"
)

type CardWorkflowTest struct {
	customerMockRepo        *mocks.CustomerRepositoryInterface
	customerTokeMockRepo    *mocks.CustomerTokenRepositoryInterface
	customerServiceMockRepo *mocks.CustomerServiceRepository
	cardMockRepo            *mocks.CardRepositoryInterface
	cardTokeMockRepo        *mocks.CardTokenRepositoryInterface
	paymentGatewayMock      *mocks.PaymentGatewayInterface
	queueMockInterface      *mocks.QueueInterface
	queueUrl                string
	app                     *bootstrap.App
	testDate                datetypes.CustomDate
	cardApplication         *mocks.ICard
}

func (c *CardWorkflowTest) initGetMockFuncs() {
	c.cardMockRepo.On("GetByID", mock.Anything, "96999c21-10ca-4391-aee9-243920ed2daf").Return(&cards.Card{
		ID:    "96999c21-10ca-4391-aee9-243920ed2daf",
		Brand: "visa",
	}, nil)
	c.cardTokeMockRepo.On("GetByCardID", mock.Anything, "96999c21-10ca-4391-aee9-243920ed2daf").Return(nil, nil)

	c.cardMockRepo.On("GetByID", mock.Anything, "cf6bc22b-83ca-42cd-b2e8-bd02d2641a36").Return(nil, &errors.ErrorOutput{Code: errors.NotFoundError, Message: "Card não encontrado", HttpStatus: http.StatusNotFound})

}

func NewCardWorkflowTest(t *testing.T) *CardWorkflowTest {
	customerMockRepo := mocks.NewCustomerRepositoryInterface(t)
	customerTokeMockRepo := mocks.NewCustomerTokenRepositoryInterface(t)
	cardMockRepo := mocks.NewCardRepositoryInterface(t)
	cardTokeMockRepo := mocks.NewCardTokenRepositoryInterface(t)
	queueMockInterface := mocks.NewQueueInterface(t)
	customerServiceMockRepo := mocks.NewCustomerServiceRepository(t)
	cardApplication := mocks.NewICard(t)
	queueUrl := "http://localhost:4566/000000000000/payment-layer-card-api-queue"

	paymentGatewayMock := mocks.NewPaymentGatewayInterface(t)
	gtws := []paymentgateway.PaymentGatewayInterface{paymentGatewayMock}
	app := bootstrap.NewApp(customerMockRepo, customerTokeMockRepo, cardMockRepo, cardTokeMockRepo, customerServiceMockRepo, gtws, queueMockInterface)

	testDate, err := datetypes.NewDateFromString("2006-01-02")
	if err != nil {
		t.Errorf("Error to create date: %v", err)
	}
	return &CardWorkflowTest{
		customerMockRepo:        customerMockRepo,
		customerTokeMockRepo:    customerTokeMockRepo,
		cardMockRepo:            cardMockRepo,
		cardTokeMockRepo:        cardTokeMockRepo,
		customerServiceMockRepo: customerServiceMockRepo,
		paymentGatewayMock:      paymentGatewayMock,
		app:                     app,
		testDate:                testDate,
		queueMockInterface:      queueMockInterface,
		queueUrl:                queueUrl,
		cardApplication:         cardApplication,
	}
}

func TestGetCardyByID(t *testing.T) {
	cardWorkflowTest := NewCardWorkflowTest(t)
	cardWorkflowTest.initGetMockFuncs()
	cc := adapters.NewCardController(cardWorkflowTest.app)
	e := echo.New()
	e.GET("/cards/:id", cc.GetCardByID)
	testTable := []struct {
		name               string
		id                 string
		expectedStatusCode int
	}{
		{
			name:               "Get card by id",
			id:                 "96999c21-10ca-4391-aee9-243920ed2daf",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Get card with invalid id",
			id:                 "cf6bc22b-83ca-42cd-b2e8-bd02d2641a36",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Get card with empty id",
			id:                 "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cards/:id", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cards/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			_ = cc.GetCardByID(c)
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rec.Result().StatusCode)
			}
		})
	}
}

func (c *CardWorkflowTest) initCreateMockFuncs() {

	mockCard := &cards.Card{
		ID:         "someId",
		CustomerID: "78539010-7ffa-4167-bd4a-eae0ad6ca55f",
	}
	customerServiceList := []*customerservice.CustomerService{
		{
			ServiceID: "789",
		},
	}

	fingerprint := mockCard.GenerateFingerprint("1234567890123456")
	mockCard.Fingerprint = fingerprint
	c.cardMockRepo.On("GetCardByFingerprint", mock.Anything, mock.AnythingOfType("string")).Return(mockCard, nil)
	c.paymentGatewayMock.On("GatewayName").Return("pagarme")
	c.queueMockInterface.On("SendMessageWithContext", mock.Anything, mock.Anything, mock.AnythingOfType("string")).Return(nil)
	c.customerServiceMockRepo.On("GetByCustomerID", mock.Anything, mock.Anything).Return(customerServiceList, nil)
	c.cardTokeMockRepo.On("GetByCardIDAndGateway", mock.Anything, "someId", "pagarme").Return(&card_token.CardToken{ID: "test"}, nil)
}

func TestCreateCard(t *testing.T) {
	cardWorkflowTest := NewCardWorkflowTest(t)
	cardWorkflowTest.initCreateMockFuncs()
	cc := adapters.NewCardController(cardWorkflowTest.app)
	e := echo.New()
	e.POST("/cards", cc.CreateCard)
	testTable := []struct {
		name               string
		input              *card_dtos.CreateCardDTOInput
		expectedStatusCode int
	}{
		{
			name:               "Create card successfully",
			expectedStatusCode: http.StatusCreated,
			input: &card_dtos.CreateCardDTOInput{
				Number:     "123456789123111111",
				CustomerID: "b289d52a-d1be-4d4c-b8cc-da89a516143b",
				Holder:     "test",
				Brand:      "Visa",
				CVV:        "123",
				ExpMonth:   "12",
				ExpYear:    "25",
				BillingAddress: &card_dtos.BillingAddress{
					State:      "MG",
					City:       "test",
					Country:    "BR",
					Line1:      "abc",
					Line2:      "def",
					District:   "test",
					PostalCode: "123456789",
				},
			},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Error to marshal json: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/cards", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = cc.CreateCard(c)

			status_code := rec.Result().StatusCode
			return_message := rec.Result().Body
			fmt.Println(return_message)
			if status_code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rec.Result().StatusCode)
			}
		})
	}
}
