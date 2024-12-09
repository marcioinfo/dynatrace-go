package card_rede_dto

import (
	"fmt"
	cardtoken_dto "payment-layer-card-api/usecases/card_token/dtos"
	"strconv"
)

func ConvertToCardTokenInputDTO(input *cardtoken_dto.CardTokenInputDTO) (*AddAndTokemizeCardInput, error) {
	var expirationMonth, expirationYear string
	monthInt, err := strconv.Atoi(input.ExpMonth)
	if err != nil {
		return nil, err
	}

	if monthInt >= 1 && monthInt <= 9 {
		expirationMonth = fmt.Sprintf("0%d", monthInt)
	} else {
		expirationMonth = fmt.Sprintf("%d", monthInt)
	}
	yearInt, err := strconv.Atoi(input.ExpYear)
	if err != nil {
		return nil, err
	}

	if yearInt >= 0 && yearInt <= 50 {
		expirationYear = fmt.Sprintf("20%02d", yearInt)
	} else if yearInt > 50 && yearInt <= 99 {
		expirationYear = fmt.Sprintf("19%02d", yearInt)
	} else {
		return nil, fmt.Errorf("Ano inválido: %s", input.ExpYear)
	}

	createCardRequest := &AddAndTokemizeCardInput{
		CreditCardNumber: input.Number,
		ExpirationMonth:  expirationMonth,
		ExpirationYear:   expirationYear,
		CustomerID:       input.CustomerToken,
		BillingName:      input.Holder,
		BillingAddress1:  input.CustomerAddress.Line1,
		BillingAddress2:  input.CustomerAddress.Line2,
		BillingCity:      input.CustomerAddress.City,
		BillingState:     input.CustomerAddress.State,
		BillingZip:       input.CustomerAddress.PostalCode,
		BillingCountry:   input.CustomerAddress.Country,
		BillingEmail:     input.CustomerEmail,
		BillingPhone:     input.CustomerPhone,
	}

	return createCardRequest, nil
}

func ConvertToZeroDollarCard(input *cardtoken_dto.CardTokenInputDTO) (*Order, error) {
	var expirationMonth, expirationYear string
	monthInt, err := strconv.Atoi(input.ExpMonth)
	if err != nil {
		return nil, err
	}

	if monthInt >= 1 && monthInt <= 9 {
		expirationMonth = fmt.Sprintf("0%d", monthInt)
	} else {
		expirationMonth = fmt.Sprintf("%d", monthInt)
	}
	yearInt, err := strconv.Atoi(input.ExpYear)
	if err != nil {
		return nil, err
	}

	if yearInt >= 0 && yearInt <= 50 {
		expirationYear = fmt.Sprintf("20%02d", yearInt)
	} else if yearInt > 50 && yearInt <= 99 {
		expirationYear = fmt.Sprintf("19%02d", yearInt)
	} else {
		return nil, fmt.Errorf("Ano inválido: %s", input.ExpYear)
	}

	return &Order{
		ZeroDollarCard: ZeroDollarCard{
			TransactionDetail: TransactionDetail{
				PayType: PayType{
					CreditCard: CreditCard{
						Number:    input.Number,
						ExpMonth:  expirationMonth,
						ExpYear:   expirationYear,
						CvvNumber: input.CVV,
					},
				},
			},
			ProcessorID:  "1",
			ReferenceNum: input.CardID,
			SaveOnFile: SaveOnFile{
				CustomerToken: input.CustomerToken,
			},
		},
	}, nil
}
