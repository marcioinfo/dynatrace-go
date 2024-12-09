package integration_utils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
)

func RequestHttp(body interface{}, path string, method string) (*http.Request, error) {
	url := os.Getenv("INTEGRATION_API_URL")

	completedPath := "/" + path
	jsonData, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url+completedPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}
	api_layer_key := os.Getenv("API_INTEGRATION_KEY")

	req.Header.Set("x-api-key", api_layer_key)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func RequestHttpWithContext(ctx context.Context, body interface{}, path string, method string) (*http.Request, error) {
	url := os.Getenv("INTEGRATION_API_URL")
	completedPath := path
	jsonData, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url+completedPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	api_layer_key := os.Getenv("API_INTEGRATION_KEY")
	req.Header.Set("x-api-key", api_layer_key)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
