package dtos

import (
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type DeleteCardTokenDTO struct {
	TokensDeleted []CardToken         `json:"tokens_deleted"` // Tokens do cartão
	Error         *errors.ErrorOutput `json:"error"`          // Erro
}
