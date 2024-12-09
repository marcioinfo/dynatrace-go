package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"go.elastic.co/apm/module/apmhttp/v2"
)

var TracingClient = apmhttp.WrapClient(http.DefaultClient)

type PayloadError struct {
	Status  string              `json:"status,omitempty"`
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
	Charges []struct {
		Status          string `json:"status,omitempty"`
		Id              string `json:"id,omitempty"`
		LastTransaction struct {
			Id              string `json:"id,omitempty"`
			Status          string `json:"status,omitempty"`
			Success         bool   `json:"success,omitempty"`
			GatewayResponse struct {
				Code   string `json:"code,omitempty"`
				Errors []struct {
					Message string `json:"message,omitempty"`
				} `json:"errors,omitempty"`
			} `json:"gateway_response,omitempty"`
			AcquirerMessage    string `json:"acquirer_message,omitempty"`
			AcquirerReturnCode string `json:"acquirer_return_code,omitempty"`
		} `json:"last_transaction,omitempty"`
	} `json:"charges,omitempty"`
}

func RequestHttp(body interface{}, path string, method string) (*http.Request, error) {
	url := os.Getenv("API_PAGARME_URL")

	completedPath := "/" + path
	jsonData, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url+completedPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	apiKey := os.Getenv("API_PAGARME_SECRET_KEY")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(apiKey+":")))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func RequestHttpWithContext(ctx context.Context, body interface{}, path string, method string) (*http.Request, error) {
	url := os.Getenv("API_PAGARME_URL")

	completedPath := "/" + path
	jsonData, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url+completedPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	apiKey := os.Getenv("API_PAGARME_SECRET_KEY")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(apiKey+":")))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func CheckForErrors(body []byte, statusCode int) *errors.ErrorOutput {
	var errorMessage string

	var data PayloadError
	err := json.Unmarshal(body, &data)
	if err != nil {
		errorMessage = "Error reading response body"
		return errors.NewError(errors.InternalServerError, err, errorMessage)
	}

	if statusCode >= 500 && statusCode <= 599 {
		errorMessage = "Internal Server Error from Pagar.me"
		return errors.NewError(errors.ExternalError, err, errorMessage)
	}

	for _, msgs := range data.Errors {
		for _, msg := range msgs {
			errorMessage = msg
			return errors.NewError(errors.BadRequest, nil, errorMessage)
		}
	}

	if data.Message != "" {
		errorMessage = data.Message
		return errors.NewError(errors.ExternalError, nil, errorMessage)
	}

	return nil
}
