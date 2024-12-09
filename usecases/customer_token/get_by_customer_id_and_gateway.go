package customer_token

import (
	"context"
	"errors"
	customer_token_entity "payment-layer-card-api/entities/customer_token"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

var (
	ErrorCustomerIdRequired    = errors.New("customer ID parâmetro é obrigatório")
	ErrorGatewayIdRequired     = errors.New("gateway parâmetro é obrigatório")
	ErrorCustomerTokenNotFound = errors.New("customer Token não encontrado")
)

type GetByCustomerIDAndGateway struct {
	CustomerTokenRepository customer_token_entity.CustomerTokenRepositoryInterface
}

func NewGetByCustomerIDAndGateway(customerTokenRepository customer_token_entity.CustomerTokenRepositoryInterface) *GetByCustomerIDAndGateway {
	return &GetByCustomerIDAndGateway{
		CustomerTokenRepository: customerTokenRepository,
	}
}

func (g *GetByCustomerIDAndGateway) Execute(ctx context.Context, customerId string, gateway string) (*customer_token_entity.CustomerToken, *layerErrors.ErrorOutput) {
	if customerId == "" {
		return nil, layerErrors.NewError(layerErrors.ParameterIsRequired, ErrorCustomerIdRequired)
	}

	if gateway == "" {
		return nil, layerErrors.NewError(layerErrors.WithoutGateway, ErrorGatewayIdRequired)
	}

	customerToken, err := g.getCustomerTokenByCustomerIDAndGateway(ctx, customerId, gateway)
	if err != nil {
		return nil, err
	}
	return customerToken, nil
}

func (g *GetByCustomerIDAndGateway) getCustomerTokenByCustomerIDAndGateway(ctx context.Context, customerId string, gateway string) (*customer_token_entity.CustomerToken, *layerErrors.ErrorOutput) {
	customerToken, err := g.CustomerTokenRepository.GetByCustomerIDAndGateway(ctx, customerId, gateway)
	if err != nil {
		return nil, err
	}
	return customerToken, nil
}
