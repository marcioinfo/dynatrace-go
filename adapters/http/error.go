package http

type ErrorWrapperInternalDTO struct {
	Message string         `json:"message"`
	Errors  []*ErrorOutput `json:"errors"`
}

type ErrorOutput struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	LogMessage string `json:"log_message"`
	HttpStatus int    `json:"http_status"`
}
