package cards

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CardRepositoryInterface interface {
	Insert(ctx context.Context, card *Card) *errors.ErrorOutput
	GetByID(ctx context.Context, id string) (*Card, *errors.ErrorOutput)
	GetByIDAndCustomerID(ctx context.Context, id string, customerId string) (*Card, *errors.ErrorOutput)
	DeleteByID(ctx context.Context, id string) (string, *errors.ErrorOutput)
	GetCardsByCustomerID(ctx context.Context, customerID string) ([]*Card, *errors.ErrorOutput)
	GetCardByFingerprint(ctx context.Context, fingerprint string) (*Card, *errors.ErrorOutput)
}
