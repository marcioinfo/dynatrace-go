package card_pagarme_dto

import (
	cardtoken_dto "payment-layer-card-api/usecases/card_token/dtos"
)

func ConvertToBillingAddress(address cardtoken_dto.HolderAddress) BillingAddress {
	return BillingAddress{
		Line1:   address.Line1,
		Line2:   address.Line2,
		ZipCode: address.PostalCode,
		City:    address.City,
		State:   address.State,
		Country: address.Country,
	}
}

func ConvertToCreditCardDTO(input *cardtoken_dto.CardTokenInputDTO) AddAndTokemizeCardInput {
	billingAddress := ConvertToBillingAddress(*input.CustomerAddress)

	return AddAndTokemizeCardInput{
		Number:          input.Number,
		HolderName:      input.Holder,
		ExpirationMonth: input.ExpMonth,
		ExpirationYear:  input.ExpYear,
		Options: Options{
			VerifyCard: true,
		},
		CVV:            input.CVV,
		BillingAddress: billingAddress,
	}
}
