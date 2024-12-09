package customer_pagarme_dto

import (
	"payment-layer-card-api/common/helpers"
	"payment-layer-card-api/entities/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

func getPhone(phone string) *Phones {
	countryCode, err := helpers.GetCountryCodeByPhoneNumber(phone)
	if err != nil {
		return &Phones{}
	}

	areaCode, err := helpers.GetAreaCodeByPhoneNumber(phone)

	if err != nil {
		return &Phones{}
	}

	phoneDetails := &Phone{
		CountryCode: countryCode,
		AreaCode:    areaCode,
		Number:      phone[4:],
	}
	return &Phones{
		HomePhone: phoneDetails,
	}
}

func ConvertToCustomerDTO(input *customer_dto.CreateCustomerDTOInput, id string) (*AddCustomerInput, *errors.ErrorOutput) {
	phone := helpers.AddPrefixToNumber(input.Phone)

	phones := getPhone(phone)

	return &AddCustomerInput{
		Name:         input.FirstName + " " + input.LastName,
		Email:        input.Email,
		Code:         id,
		Document:     input.Document,
		DocumentType: helpers.CheckDocumentType(input.Document),
		Type:         "individual",
		Gender:       input.Gender,
		Address: Address{
			Country: input.Address.Country,
			State:   input.Address.State,
			City:    input.Address.City,
			ZipCode: string(input.Address.PostalCode),
			Line1:   input.Address.Line1,
			Line2:   input.Address.Line2,
		},
		Phones:    *phones,
		Birthdate: input.BirthDate.Format("01/02/2006"),
	}, nil
}

func ConvertToEditCustomerDTO(customerInput *customers.Customer) *EditCustomerDTO {
	phone := helpers.AddPrefixToNumber(customerInput.Phone)

	phones := getPhone(phone)
	return &EditCustomerDTO{
		Name:      customerInput.Name,
		Phones:    *phones,
		Email:     customerInput.Email,
		Gender:    customerInput.Gender,
		BirthDate: customerInput.BirthDate.Format("01/02/2006"),
		Document:  customerInput.Document,
		Type:      "individual", //TODO: mudar isso quando aceitarmos clientes pj, dai teremos que salvar o type.
	}
}
