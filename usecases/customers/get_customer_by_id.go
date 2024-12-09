package customers_usecase

import (
	"context"
	"errors"

	"payment-layer-card-api/entities/customers"

	"github.com/adhfoundation/layer-tools/log"
	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

var (
	ErrorIdRequired       = errors.New("id é obrigatório")
	ErrorCustomerNotFound = errors.New("customer não encontrado")
)

type GetCustomerByID struct {
	CustomerRepository customers.CustomerRepositoryInterface
}

func NewGetCustomerByID(customerRepository customers.CustomerRepositoryInterface) *GetCustomerByID {
	return &GetCustomerByID{
		CustomerRepository: customerRepository,
	}
}

func (gc *GetCustomerByID) Execute(ctx context.Context, customerID string) (*customers.Customer, *layerErrors.ErrorOutput) {
	if customerID == "" {
		err := layerErrors.NewError(layerErrors.ParameterIsRequired, nil, "CustomerID parâmetro é obrigatório.")
		log.Error(ctx, err.LogMessageToError()).Msgf("Erro ao consultar cliente pelo CustomerID vazio")
		return nil, err
	}

	customer, err := gc.getCustomerInDatabase(ctx, customerID)
	if err != nil {
		if err.Code == layerErrors.Unauthorized {
			log.Info(ctx).Msgf("O cliente não foi encontrao para o CustomerID (%v)", customerID)
			return nil, err
		}
		log.Error(ctx, err.LogMessageToError()).Msgf("Erro ao consultar cliente pelo CustomerID (%v), erro: (%v)", customerID, err)
		return nil, err
	}

	return customer, nil
}

func (gc *GetCustomerByID) getCustomerInDatabase(ctx context.Context, id string) (*customers.Customer, *layerErrors.ErrorOutput) {
	customer, err := gc.CustomerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
