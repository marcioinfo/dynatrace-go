package customers

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CustomerRepositoryInterface interface {
	Insert(ctx context.Context, customer *Customer) *errors.ErrorOutput
	GetByID(ctx context.Context, id string) (*Customer, *errors.ErrorOutput)
	GetCustomerByDocument(ctx context.Context, document string, serviceID string) (*Customer, *errors.ErrorOutput)
	GetCustomerByEmail(ctx context.Context, email string) (*Customer, *errors.ErrorOutput)
	Update(ctx context.Context, customer *Customer) *errors.ErrorOutput
}
