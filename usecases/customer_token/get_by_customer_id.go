package customer_token

import (
	"context"
	customertoken "payment-layer-card-api/entities/customer_token"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type GetByCustomerID struct {
	CustomerTokenRepository customertoken.CustomerTokenRepositoryInterface
}

func NewGetByCustomerID(customerTokenRepository customertoken.CustomerTokenRepositoryInterface) *GetByCustomerID {
	return &GetByCustomerID{
		CustomerTokenRepository: customerTokenRepository,
	}
}

func (g *GetByCustomerID) Execute(ctx context.Context, customerID string) ([]*customertoken.CustomerToken, *layerErrors.ErrorOutput) {
	if customerID == "" {
		return nil, layerErrors.NewError(layerErrors.ParameterIsRequired, ErrorCustomerIdRequired)
	}

	customerToken, err := g.getCustomerTokenByCustomerIDInDatabase(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return customerToken, nil
}

func (g *GetByCustomerID) getCustomerTokenByCustomerIDInDatabase(ctx context.Context, customerID string) ([]*customertoken.CustomerToken, *layerErrors.ErrorOutput) {
	customerToken, err := g.CustomerTokenRepository.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return customerToken, nil
}
