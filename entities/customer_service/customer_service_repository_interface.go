package customerservice

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CustomerServiceRepository interface {
	Insert(ctx context.Context, customerService *CustomerService) *errors.ErrorOutput
	GetByCustomerID(ctx context.Context, customerID string) ([]*CustomerService, *errors.ErrorOutput)
	GetByServiceID(ctx context.Context, serviceID string) ([]*CustomerService, *errors.ErrorOutput)
	Update(ctx context.Context, customerService *CustomerService) *errors.ErrorOutput
	Delete(ctx context.Context, customerService *CustomerService) *errors.ErrorOutput
	GetByCustomerAndServiceID(ctx context.Context, customerID, serviceID string) (*CustomerService, *errors.ErrorOutput)
}
