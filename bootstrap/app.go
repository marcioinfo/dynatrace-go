package bootstrap

import (
	"os"
	cardtoken "payment-layer-card-api/entities/card_token"
	"payment-layer-card-api/entities/cards"
	customerservice "payment-layer-card-api/entities/customer_service"
	customertoken "payment-layer-card-api/entities/customer_token"
	"payment-layer-card-api/entities/customers"
	paymentgateway "payment-layer-card-api/entities/payment_gateway"
	"payment-layer-card-api/entities/queue"
	"strconv"

	"github.com/adhfoundation/layer-tools/apmtracer"
	"github.com/go-playground/validator/v10"
	"go.elastic.co/apm/v2"
)

type App struct {
	Validator           *validator.Validate
	CustomerRepo        customers.CustomerRepositoryInterface
	CustomerTokenRepo   customertoken.CustomerTokenRepositoryInterface
	CardRepo            cards.CardRepositoryInterface
	CardTokenRepo       cardtoken.CardTokenRepositoryInterface
	CustomerServiceRepo customerservice.CustomerServiceRepository
	PaymentGateways     []paymentgateway.PaymentGatewayInterface
	QueueService        queue.QueueInterface
	ApmTracer           *apm.Tracer
}

func NewApp(customerRepo customers.CustomerRepositoryInterface,
	customerTokenRepo customertoken.CustomerTokenRepositoryInterface,
	cardRepo cards.CardRepositoryInterface,
	cartTokenRepo cardtoken.CardTokenRepositoryInterface,
	customerServiceRepo customerservice.CustomerServiceRepository,
	paymentGateways []paymentgateway.PaymentGatewayInterface,
	queueService queue.QueueInterface,
) *App {

	return &App{
		Validator:           buildValidator(),
		CustomerRepo:        customerRepo,
		CustomerTokenRepo:   customerTokenRepo,
		CardRepo:            cardRepo,
		CardTokenRepo:       cartTokenRepo,
		CustomerServiceRepo: customerServiceRepo,
		PaymentGateways:     paymentGateways,
		QueueService:        queueService,
		ApmTracer:           apmtracer.Tracer,
	}
}

func (app *App) GetCustomerRetryCount() int {
	customerRetryCount, _ := strconv.Atoi(os.Getenv("CUSTOMER_RETRY_COUNT"))
	if customerRetryCount == 0 {
		customerRetryCount = 3
	}
	return customerRetryCount
}

func (app *App) GetCustomerRetryDelayInMilliseconds() int {
	customerRetryDelay, _ := strconv.Atoi(os.Getenv("CUSTOMER_RETRY_DELAY"))
	return customerRetryDelay
}

func (app *App) GetCardRetryCount() int {
	cardRetryCount, _ := strconv.Atoi(os.Getenv("CARD_RETRY_COUNT"))
	return cardRetryCount
}

func (app *App) GetCardRetryDelayInMilliseconds() int {
	cardRetryDelay, _ := strconv.Atoi(os.Getenv("CARD_RETRY_DELAY"))
	return cardRetryDelay
}

func (app *App) GetSyncCardQueueURL() string {
	queueURL := os.Getenv("CREATE_CARD_INTEGRATION_QUEUE_URL")
	return queueURL
}
