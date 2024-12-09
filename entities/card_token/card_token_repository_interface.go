package card_token

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CardTokenRepositoryInterface interface {
	Insert(ctx context.Context, cardToken *CardToken) *errors.ErrorOutput
	GetByCardID(ctx context.Context, cardID string) ([]*CardToken, *errors.ErrorOutput)
	DeleteByCardToken(ctx context.Context, cardID string) *errors.ErrorOutput
	GetByCardIDAndGateway(ctx context.Context, cardID string, gateway string) (*CardToken, *errors.ErrorOutput)
}
