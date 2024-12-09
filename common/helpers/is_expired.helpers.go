package helpers

import (
	"strconv"
	"time"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

func IsExpired(month, year string) (bool, *errors.ErrorOutput) {
	expMonth, err := strconv.Atoi(month)
	if err != nil {
		return false, errors.NewError(errors.ValidationEntityError, err, "erro convertendo mês")

	}

	expYear, err := strconv.Atoi(year)
	if err != nil {
		return false, errors.NewError(errors.ValidationEntityError, err, "erro convertendo ano")
	}

	expYear += 2000

	currentTime := time.Now()
	currentMonth := int(currentTime.Month())
	currentYear := currentTime.Year()

	if expYear < currentYear {
		return true, nil
	}

	if expYear == currentYear && expMonth < currentMonth {
		return true, nil
	}

	return false, nil
}
