package customer_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	adapters "payment-layer-card-api/adapters/http"
	"payment-layer-card-api/bootstrap"
	customerservice "payment-layer-card-api/entities/customer_service"
	"payment-layer-card-api/entities/customers"
	paymentgateway "payment-layer-card-api/entities/payment_gateway"
	"payment-layer-card-api/mocks"
	customertoken_dto "payment-layer-card-api/usecases/customer_token/dtos"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"strings"
	"testing"
	"time"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type CustomerWorkflowTest struct {
	customerMockRepo        *mocks.CustomerRepositoryInterface
	customerTokeMockRepo    *mocks.CustomerTokenRepositoryInterface
	cardMockRepo            *mocks.CardRepositoryInterface
	cardTokeMockRepo        *mocks.CardTokenRepositoryInterface
	customerServiceMockRepo *mocks.CustomerServiceRepository
	paymentGatewayMock      *mocks.PaymentGatewayInterface
	queueService            *mocks.QueueInterface
	app                     *bootstrap.App
	testDate                datetypes.CustomDate
}

func (c *CustomerWorkflowTest) initCreateMockFuncs() {
	c.customerMockRepo.On("GetCustomerByEmail", mock.Anything, "john@example.com").Return(nil, nil)
	c.customerMockRepo.On("Insert", mock.Anything, mock.AnythingOfType("*customers.Customer")).Return(nil)
	c.customerTokeMockRepo.On("GetByCustomerIDAndGateway", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, nil)
	c.customerTokeMockRepo.On("Insert", mock.Anything, mock.AnythingOfType("*customer_token.CustomerToken")).Return(nil)
	c.customerTokeMockRepo.On("Insert", mock.Anything, mock.AnythingOfType("*customer_token.CustomerToken")).Return(layerErrors.NewPaymentLayerError(layerErrors.InternalServerError, "Internal error"))
	c.paymentGatewayMock.On("GatewayName").Return("mock")
	c.paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.AnythingOfType("*customer_dto.CreateCustomerDTOInput"), mock.AnythingOfType("string")).Return(&customertoken_dto.CustomerTokenOutputDTO{
		Token: "sample_token",
	}, 201, (*errors.ErrorOutput)(nil))
	c.paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.AnythingOfType("*customer_dto.CreateCustomerDTOInput"), mock.AnythingOfType("string")).Return(nil, 400, layerErrors.NewPaymentLayerError(layerErrors.InvalidCustomer, "Invalid customer"))
	c.queueService.On("SendMessageWithContext", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	c.customerMockRepo.On("GetByID", mock.Anything, mock.Anything).Return(&customers.Customer{
		ID:       "8fd18102-c0a4-4c64-9a4e-dbf06ffd6e11",
		Name:     "John Doe",
		Document: "12345678901",
	}, nil)
	c.customerServiceMockRepo.On("GetByCustomerAndServiceID", mock.Anything, mock.Anything, mock.Anything).Return(&customerservice.CustomerService{
		ID:         "cs-id-1",
		ServiceID:  "service-id-1",
		CustomerID: "customer-id-1",
		Name:       "John Doe",
		Document:   "12345678901",
		Email:      "teste@teste.com",
		BirthDate:  datetypes.CustomDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
		CreatedAt:  datetypes.CustomDateTime(time.Now()),
		UpdatedAt:  datetypes.CustomDateTime(time.Now()),
		Gender:     "male",
	}, nil)
	c.customerServiceMockRepo.On("Insert", mock.Anything, mock.Anything).Return(nil)
	c.customerTokeMockRepo.On("GetByCustomerID", mock.Anything, mock.Anything).Return(nil, nil)
}

func (c *CustomerWorkflowTest) initGetMockFuncs() {
	c.customerMockRepo.On("GetByID", mock.Anything, "96999c21-10ca-4391-aee9-243920ed2daf").Return(&customers.Customer{
		ID:       "96999c21-10ca-4391-aee9-243920ed2daf",
		Name:     "John Doe",
		Document: "12345678901",
	}, nil)
	c.customerMockRepo.On("GetByID", mock.Anything, "cf6bc22b-83ca-42cd-b2e8-bd02d2641a36").Return(nil, nil)
	c.customerTokeMockRepo.On("GetByCustomerID", mock.Anything, mock.Anything).Return(nil, nil)
}

func NewCustomerWorkflowTest(t *testing.T) *CustomerWorkflowTest {
	customerMockRepo := mocks.NewCustomerRepositoryInterface(t)
	customerTokeMockRepo := mocks.NewCustomerTokenRepositoryInterface(t)
	cardMockRepo := mocks.NewCardRepositoryInterface(t)
	cardTokeMockRepo := mocks.NewCardTokenRepositoryInterface(t)
	paymentGatewayMock := mocks.NewPaymentGatewayInterface(t)
	customerServiceMockRepo := mocks.NewCustomerServiceRepository(t)
	queueMock := mocks.NewQueueInterface(t)
	gtws := []paymentgateway.PaymentGatewayInterface{paymentGatewayMock}
	app := bootstrap.NewApp(customerMockRepo, customerTokeMockRepo, cardMockRepo, cardTokeMockRepo, customerServiceMockRepo, gtws, queueMock)

	testDate, err := datetypes.NewDateFromString("2006-01-02")
	if err != nil {
		t.Errorf("Error to create date: %v", err)
	}
	return &CustomerWorkflowTest{
		customerMockRepo:        customerMockRepo,
		customerTokeMockRepo:    customerTokeMockRepo,
		cardMockRepo:            cardMockRepo,
		cardTokeMockRepo:        cardTokeMockRepo,
		customerServiceMockRepo: customerServiceMockRepo,
		paymentGatewayMock:      paymentGatewayMock,
		app:                     app,
		queueService:            queueMock,
		testDate:                testDate,
	}
}

func TestGetCustomerByID(t *testing.T) {
	cardAPITest := NewCustomerWorkflowTest(t)
	cardAPITest.initGetMockFuncs()
	cc := adapters.NewCustomerController(cardAPITest.app)
	e := echo.New()
	e.GET("/customers/:id", cc.GetCustomerByID)
	testTable := []struct {
		name               string
		id                 string
		expectedStatusCode int
		getCalledTimes     int
	}{
		{
			name:               "Get customer by id",
			id:                 "96999c21-10ca-4391-aee9-243920ed2daf",
			expectedStatusCode: http.StatusOK,
			getCalledTimes:     1,
		},
		{
			name:               "Get customer invalid id",
			id:                 "cf6bc22b-83ca-42cd-b2e8-bd02d2641a36",
			expectedStatusCode: http.StatusNotFound,
			getCalledTimes:     404,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/customer/:id", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/customers/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			_ = cc.GetCustomerByID(c)
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rec.Result().StatusCode)
			}
		})
	}
}

func TestCreateCustomer(t *testing.T) {

	cardAPITest := NewCustomerWorkflowTest(t)
	cardAPITest.initCreateMockFuncs()
	cc := adapters.NewCustomerController(cardAPITest.app)
	e := echo.New()
	e.POST("/customers", cc.CreateCustomer)
	testTable := []struct {
		name               string
		input              customer_dto.CreateCustomerDTOInput
		wantErr            bool
		expectedStatusCode int
		validateBody       string
	}{
		{
			name: "Create customer with success",
			input: customer_dto.CreateCustomerDTOInput{
				FirstName: "John",
				LastName:  "Doe",
				Document:  "93998295020",
				BirthDate: cardAPITest.testDate,
				Gender:    "male",
				ServiceID: "b289d52a-d1be-4d4c-b8cc-da89a516143b",
				Phone:     "5511999999999",
				Email:     "john@example.com",
				Address: customer_dto.Address{
					State:      "SP",
					City:       "São Paulo",
					Country:    "BR",
					Line1:      "Rua João Silva, 123",
					Line2:      "Apto 123",
					District:   "Centro",
					PostalCode: "12345678",
				},
			},
			wantErr:            false,
			expectedStatusCode: 201,
			validateBody:       "",
		},
		{
			name: "Create customer with invalid document",
			input: customer_dto.CreateCustomerDTOInput{
				FirstName: "John",
				LastName:  "Doe",
				Document:  "123",
				ServiceID: "b289d52a-d1be-4d4c-b8cc-da89a516143b",
				BirthDate: cardAPITest.testDate,
				Gender:    "male",
				Phone:     "5511999999999",
				Email:     "johnd@example.com",
				Address: customer_dto.Address{
					State:      "SP",
					City:       "São Paulo",
					Country:    "BR",
					Line1:      "Rua João Silva, 123",
					Line2:      "Apto 123",
					District:   "Centro",
					PostalCode: "12345678",
				},
			},
			wantErr:            true,
			expectedStatusCode: 400,
			validateBody:       "",
		},
		{
			name: "Create customer with invalid email",
			input: customer_dto.CreateCustomerDTOInput{
				FirstName: "John",
				LastName:  "Doe",
				Document:  "93998295020",
				ServiceID: "b289d52a-d1be-4d4c-b8cc-da89a516143b",
				BirthDate: cardAPITest.testDate,
				Gender:    "male",
				Phone:     "5511999999999",
				Email:     "johnexample.com",
				Address: customer_dto.Address{
					State:      "SP",
					City:       "São Paulo",
					Country:    "BR",
					Line1:      "Rua João Silva, 123",
					Line2:      "Apto 123",
					District:   "Centro",
					PostalCode: "12345678",
				},
			},
			wantErr:            true,
			expectedStatusCode: 400,
			validateBody:       "",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Error to marshal input: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/customers")
			_ = cc.CreateCustomer(c)
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rec.Result().StatusCode)
			}
		})
	}

}
