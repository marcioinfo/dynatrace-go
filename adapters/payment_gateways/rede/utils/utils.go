package utils

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"payment-layer-card-api/adapters/payment_gateways/rede/templates"
	"strings"
	"text/template"
	"time"

	"go.elastic.co/apm/module/apmhttp/v2"
)

type apiRequest struct {
	Command string
	Request string
}

type transactionRequest struct {
	Command *string
	Request string
}

var TracingClient = apmhttp.WrapClient(http.DefaultClient)

func BuildApiRequestXmlBody(body string, command string) string {
	template, err := template.New("requestTemplate").Parse(templates.GetApiRequestTemplateValue())
	if err != nil {
		panic(err)
	}
	data := apiRequest{
		Command: command,
		Request: string(body),
	}

	var buf strings.Builder
	err = template.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	xmlRequest := buf.String()

	return xmlRequest
}

func BuildTransactionXmlBody(body string) string {
	template, err := template.New("requestTemplate").Parse(templates.GetTransactionRequestTemplateValue())
	if err != nil {
		panic(err)
	}
	data := transactionRequest{
		Request: string(body),
	}

	var buf strings.Builder
	err = template.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	xmlRequest := buf.String()

	return xmlRequest
}

func RequestHttp(body string, path string) (*http.Request, error) {
	url := os.Getenv("API_REDE_URL")
	req, err := http.NewRequest(http.MethodPost, url+path, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	return req, nil
}

func RequestHttpWithContext(ctx context.Context, body string, path string) (*http.Request, error) {

	url := os.Getenv("API_REDE_URL")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+path, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	return req, nil
}

func ConvertDateToRedeFormat(date string) (string, error) {
	parsedDate, _ := time.Parse("02/01/2006", date)
	formattedDate := parsedDate.Format("01/02/2006")

	return formattedDate, nil
}

func ConvertGenderToRedeFormat(gender string) string {
	switch gender {
	case "male":
		return "m"
	case "female":
		return "f"
	default:
		return "f"
	}
}
