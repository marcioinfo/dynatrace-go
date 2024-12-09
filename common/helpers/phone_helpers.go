package helpers

import (
	"errors"
	"payment-layer-card-api/common/helpers/prefixes"
	"regexp"
)

func GetCountryCodeByPhoneNumber(phoneNumber string) (string, error) {
	if len(phoneNumber) > 2 {
		return phoneNumber[0:2], nil
	}

	return "", errors.New("número de telefone inválido")
}

func GetAreaCodeByPhoneNumber(phoneNumber string) (string, error) {
	if len(phoneNumber) > 4 {
		return phoneNumber[2:4], nil
	}

	return "", errors.New("número de telefone inválido")
}

func IsValidPhone(phoneNumber string) bool {
	regex := regexp.MustCompile(`[^\p{N} ]+$`)

	clearNumber := regex.ReplaceAllString(phoneNumber, "")

	return len(clearNumber) >= 12
}

func AddPrefixToNumber(phone string) string {
	if len(phone) > 0 && phone[:2] != "55" {
		phone = prefixes.BR + phone
	}

	return phone
}
