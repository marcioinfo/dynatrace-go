package http

import (
	"context"
	"net/http"
	"os"
	"payment-layer-card-api/bootstrap"
	"payment-layer-card-api/common/helpers"
	_ "payment-layer-card-api/docs"
	"strconv"
	"strings"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/layer-tools/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.elastic.co/apm/module/apmechov4/v2"
)

type Server struct {
	app *bootstrap.App
	ech *echo.Echo
}

func NewServer(app *bootstrap.App) *Server {
	log.Info(context.Background()).Msg("Iniciando HTTP server...")
	return &Server{
		app: app,
		ech: echo.New(),
	}
}

// @title Card Layer API
// @version 1.0
// @description Layer para integrações de pagamento Afya.

// @contact.name Foundation Squad

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

//@securityDefinitions.apikey x-api-key
//@in header
//@name x-api-key
//@Security x-api-key

// @BasePath /
// @schemes https http
func (s *Server) Start() error {
	port := os.Getenv("HTTP_PORT")

	if port == "" {
		port = "8085"
	}

	enableHealhtzLogging := os.Getenv("ENABLE_HEALTHZ_LOGGING")

	s.ech.Use(middlewares.LogMiddlewareRequestLogger(enableHealhtzLogging, helpers.DISABLE_LOG_ROUTE))

	s.ech.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	s.ech.Use(apmechov4.Middleware(apmechov4.WithTracer(s.app.ApmTracer)))
	s.ech.Use(middlewares.AuditMiddleware())
	s.ech.Use(middlewares.ApmMiddlewareCaptureBody())
	s.ech.Use(middlewares.PanicRecoveryMiddleware())

	var appliedRateLimit bool
	rateLimitEnabled := os.Getenv("RATE_LIMIT_ENABLED")
	if rateLimitEnabled == "enabled" {
		appliedRateLimit = true
	} else {
		appliedRateLimit = false
	}

	allowedIPs := strings.Split(os.Getenv("RATE_LIMIT_EXEMPT_IPS"), ",")
	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		rateLimit = 10
	}

	customerController := NewCustomerController(s.app)
	cardController := NewCardController(s.app)
	customerServiceController := NewCustomerServiceController(s.app)
	s.ech.POST(
		"/customers",
		customerController.CreateCustomer,
		middlewares.PermissionCheckMiddleware("card_api:customers:create", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	)
	s.ech.GET(
		"/customers/:id",
		customerController.GetCustomerByID,
		middlewares.PermissionCheckMiddleware("card_api:customers:getById", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	)
	s.ech.GET(
		"/customers/document/:document/serviceID/:serviceID",
		customerController.GetCustomerByDocument,
		middlewares.PermissionCheckMiddleware("card_api:customers:getByDocument", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	)
	s.ech.DELETE(
		"/customers/:customerId/cards/:cardId",
		cardController.DeleteCardByIDAndCustomerID,
		middlewares.PermissionCheckMiddleware("card_api:customers:delete", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	)
	// s.ech.PUT(
	// 	"/customers/:id",
	// 	customerController.UpdateCustomer,
	// 	middlewares.PermissionCheckMiddleware("card_api:customers:update", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	// )

	if appliedRateLimit {
		s.ech.POST(
			"/cards",
			middlewares.RateLimitMiddleware(cardController.CreateCard, allowedIPs, rateLimit),
			middlewares.PermissionCheckMiddleware("card_api:cards:create", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
		)
	} else {
		s.ech.POST(
			"/cards",
			cardController.CreateCard,
			middlewares.PermissionCheckMiddleware("card_api:cards:create", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
		)
	}

	s.ech.GET(
		"/cards/:id",
		cardController.GetCardByID,
		middlewares.PermissionCheckMiddleware("card_api:cards:getById", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	)
	s.ech.GET(
		"/cards/:cardId/customer/:customerId",
		cardController.GetCardByIDAndCustomerID,
		middlewares.PermissionCheckMiddleware("card_api:cards:getByIdAndCustomerId", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	)

	s.ech.GET(
		"/customer-service/customer/:customerId/service/:serviceId",
		customerServiceController.GetCustomerServiceByServiceIdAndCustomerId,
		middlewares.PermissionCheckMiddleware("card_api:customer-service:getByServiceIdAndCustomerId", os.Getenv("AUTH_BRIDGE_URL"), os.Getenv("API_INTEGRATION_KEY")),
	)

	s.ech.GET("/healthz", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})

	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		s.ech.GET("/docs/*", echoSwagger.WrapHandler)
	}
	return s.ech.Start(":" + port)
}
