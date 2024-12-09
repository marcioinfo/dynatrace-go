package customer_rede_dto

import (
	"payment-layer-card-api/entities/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"strings"
)

func ConvertToCreateCustomerDTO(input *customer_dto.CreateCustomerDTOInput, id string) *AddConsumerInput {

	var sex string

	if input.Gender != "" {
		sex = strings.ToUpper(string(input.Gender[0]))
	}

	return &AddConsumerInput{
		CustomerIDExt: id,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Address1:      input.Address.Line1,
		Address2:      input.Address.Line2,
		City:          input.Address.City,
		State:         input.Address.State,
		Zip:           input.Address.PostalCode,
		Country:       input.Address.Country,
		Phone:         input.Phone,
		Email:         input.Email,
		Dob:           input.BirthDate.Format("01/02/2006"),
		Sex:           sex,
	}
}

func ConvertToUpdateCustomerDTO(input *customers.Customer, token string) *UpdateCustomerDTO {
	fistName := strings.Split(input.Name, " ")[0]
	lastName := strings.Split(input.Name, " ")[1]

	var sex string

	if input.Gender != "" {
		sex = strings.ToUpper(string(input.Gender[0]))
	}

	return &UpdateCustomerDTO{
		CustomerID: token,
		FirstName:  fistName,
		LastName:   lastName,
		Phone:      input.Phone,
		Email:      input.Email,
		Dob:        input.BirthDate.Format("01/02/2006"),
		Sex:        sex,
	}
}
