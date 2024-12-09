package factories

import (
	"database/sql"
	"fmt"
	"os"
	"payment-layer-card-api/adapters/aws"
	postgres "payment-layer-card-api/adapters/db"
	"payment-layer-card-api/adapters/http"
	"payment-layer-card-api/adapters/payment_gateways/pagarme"
	"payment-layer-card-api/adapters/payment_gateways/rede"
	"payment-layer-card-api/adapters/queue/sqs"
	"payment-layer-card-api/bootstrap"
	paymentgateway "payment-layer-card-api/entities/payment_gateway"
)

func NewServerFactory(
	db *sql.DB,
) *http.Server {
	awsConnection, err := aws.NewAWSConnection()
	if err != nil {
		fmt.Printf("erro: %+v", err.Error())
	}
	sqs := sqs.NewSqsAWS(awsConnection)
	pagarmeURL := os.Getenv("API_PAGARME_URL")
	pagarmeSecretKey := os.Getenv("API_PAGARME_SECRET_KEY")
	pagarmeAuthorization := os.Getenv("API_PAGARME_AUTHORIZATION")
	redeUrl := os.Getenv("API_REDE_URL")
	redeMerchantKey := os.Getenv("API_REDE_MERCHANT_KEY")
	redeMerchantId := os.Getenv("API_REDE_MERCHANT_ID")
	customerRepo := postgres.NewCustomerRepository(db)
	customerTokeRepo := postgres.NewCustomerTokenRepository(db)
	customerServiceRepo := postgres.NewCustomerServiceRepository(db)
	cardRepo := postgres.NewCardRepository(db)
	cardTokenRepo := postgres.NewCardTokenRepository(db)
	paymentGateways := []paymentgateway.PaymentGatewayInterface{
		rede.NewRedeAdapter(redeUrl, redeMerchantKey, redeMerchantId),
		pagarme.NewPagarmeAdapter(pagarmeURL, pagarmeSecretKey, pagarmeAuthorization),
	}

	app := bootstrap.NewApp(
		customerRepo,
		customerTokeRepo,
		cardRepo,
		cardTokenRepo,
		customerServiceRepo,
		paymentGateways,
		sqs,
	)
	httpServer := http.NewServer(app)
	return httpServer
}
