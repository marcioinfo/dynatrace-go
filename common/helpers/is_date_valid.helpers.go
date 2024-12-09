package helpers

import (
	"fmt"
	"strconv"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

func IsValidDate(month, year string) (bool, *errors.ErrorOutput) {
	expMonth, err := strconv.Atoi(month)
	if err != nil {
		return false, errors.NewError(errors.ValidationEntityError, err, fmt.Sprintf("mês inválido %v", err))
	}

	expYear, err := strconv.Atoi(year)
	if err != nil {
		return false, errors.NewError(errors.ValidationEntityError, err, fmt.Sprintf("ano inválido: %v", err))
	}

	if expMonth < 1 || expMonth > 12 {
		return false, errors.NewError(errors.ValidationEntityError, nil, "mês deve estar entre 1 e 12")
	}

	if expYear < 0 || expYear > 99 {
		return false, errors.NewError(errors.ValidationEntityError, nil, "ano deve ter dois dígitos")
	}

	return true, nil
}
