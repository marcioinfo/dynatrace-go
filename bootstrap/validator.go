package bootstrap

import (
	"context"
	"payment-layer-card-api/common/helpers"
	"payment-layer-card-api/common/validators"
	"regexp"
	"time"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/go-playground/validator/v10"
	"github.com/klassmann/cpfcnpj"
)

func buildValidator() *validator.Validate {
	validate := validator.New()

	err := validate.RegisterValidation("cep", func(fl validator.FieldLevel) bool {
		value := fl.Field()
		match, _ := regexp.MatchString("\\d{8}", value.String())
		return match
	})
	if err != nil {
		log.Error(context.Background(), err).Msg("Error ao registrar cep validator" + err.Error())
	}

	err = validate.RegisterValidation("states", func(fl validator.FieldLevel) bool {
		value := fl.Field()
		return validators.IsState(value.String())
	})
	if err != nil {
		log.Error(context.Background(), err).Msg("Error ao registrar states validator" + err.Error())
	}

	err = validate.RegisterValidation("currentYear", func(fl validator.FieldLevel) bool {
		t := time.Now()
		year := t.Year()

		value := fl.Field()

		return int(value.Int()) >= year
	})
	if err != nil {
		log.Error(context.Background(), err).Msg("Error ao registrar currentYear validator" + err.Error())
	}

	err = validate.RegisterValidation("cpf_or_cnpj", func(fl validator.FieldLevel) bool {
		value := fl.Field()
		cpf := cpfcnpj.NewCPF(value.String())
		cnpj := cpfcnpj.NewCNPJ(value.String())
		return cpf.IsValid() || cnpj.IsValid()
	})
	if err != nil {
		log.Error(context.Background(), err).Msg("Error ao registrar cpf_or_cnpj validator" + err.Error())
	}

	err = validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		value := fl.Field()
		return helpers.IsValidPhone(value.String())

	})
	if err != nil {
		log.Error(context.Background(), err).Msg("Error ao registrar phone validator" + err.Error())
	}

	err = validate.RegisterValidation("year", func(fl validator.FieldLevel) bool {
		value := fl.Field()

		regex := regexp.MustCompile(`[^\p{N} ]+$`)

		clearNumber := regex.ReplaceAllString(value.String(), "")

		return len(clearNumber) == 4
	})
	if err != nil {
		log.Error(context.Background(), err).Msg("Error ao registrar year validator" + err.Error())
	}

	err = validate.RegisterValidation("month", func(fl validator.FieldLevel) bool {
		value := fl.Field()

		regex := regexp.MustCompile(`[^\p{N} ]+$`)

		clearNumber := regex.ReplaceAllString(value.String(), "")

		return len(clearNumber) == 2
	})
	if err != nil {
		log.Error(context.Background(), err).Msg("Error ao registrar month validator" + err.Error())
	}

	return validate
}
