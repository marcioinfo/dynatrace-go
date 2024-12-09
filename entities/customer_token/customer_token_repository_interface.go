package customer_token

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CustomerTokenRepositoryInterface interface {
	Insert(ctx context.Context, customerToken *CustomerToken) *errors.ErrorOutput
	GetByCustomerID(ctx context.Context, customerID string) ([]*CustomerToken, *errors.ErrorOutput)
	GetByCustomerIDAndGateway(ctx context.Context, customerID string, gateway string) (*CustomerToken, *errors.ErrorOutput)
}
