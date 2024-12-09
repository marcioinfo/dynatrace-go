package integrationapi

import (
	"context"
	"io"
	"net/http"
	integration_utils "payment-layer-card-api/adapters/integration_api/utils"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type IntegrationApiAdapter struct {
	url string
}

func NewIntegrationApiAdapter(url string) *IntegrationApiAdapter {
	return &IntegrationApiAdapter{
		url: url,
	}
}

func (i *IntegrationApiAdapter) DeleteCardInIntegration(ctx context.Context, cardID string, customerID string) (*http.Response, error) {
	query := "/customers/" + customerID + "/cards/" + cardID
	req, err := integration_utils.RequestHttpWithContext(ctx, nil, query, http.MethodDelete)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewError(errors.ExternalError, err, "Erro ao deletar cartão no integration api")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return nil, errors.NewError(errors.ExternalError, nil, "Erro ao deletar cartão no integration api: "+string(body))
	}
	return resp, nil
}
