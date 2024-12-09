package dtos

import "strings"

type ResponseErrorDTO struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func GetAllErrorsByResponse(responseError *ResponseErrorDTO) string {
	var allErrors []string
	for _, errorList := range responseError.Errors {
		allErrors = append(allErrors, errorList...)
	}
	return strings.Join(allErrors, ",\n")
}
