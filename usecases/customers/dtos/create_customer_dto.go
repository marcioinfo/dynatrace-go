package customer_dto

import (
	"github.com/adhfoundation/layer-tools/datetypes"
)

type CreateCustomerDTOInput struct {
	FirstName string               `json:"first_name" validate:"required,max=50" example:"João"`                                // Nome do cliente
	LastName  string               `json:"last_name" validate:"required,max=50" example:"Silva"`                                // Sobrenome do cliente
	Document  string               `json:"document" validate:"required,cpf_or_cnpj,len=11|len=14" example:"12345678901"`        // CPF ou CNPJ
	BirthDate datetypes.CustomDate `json:"birth_date" validate:"required" example:"2000-11-02"`                                 // Data de nascimento
	Gender    string               `json:"gender" validate:"oneof=male female ''" example:"male || female"`                     // Gênero
	Phone     string               `json:"phone" validate:"required,numeric,gte=10,max=13" example:"11999999999"`               // Telefone
	Email     string               `json:"email" validate:"required,email,max=100" example:"teste@gmail.com"`                   // E-mail
	Address   Address              `json:"address" validate:"required"`                                                         // Endereço do cliente
	ServiceID string               `json:"service_id" validate:"required,uuid4" example:"835524dd-d897-45ed-a68e-2e1a8bf522d7"` // ID do serviço
}

type Address struct {
	State         string `json:"state" validate:"required,states,len=2" example:"SP"`                 // Estado do endereço
	City          string `json:"city" validate:"required,max=50" example:"São Paulo"`                 // Cidade do endereço
	Country       string `json:"country" validate:"required,iso3166_1_alpha2,len=2" example:"BR"`     // País do endereço
	Line1         string `json:"line_1" validate:"required,max=100" example:"Rua João Silva, 123"`    // Linha 1 do endereço
	Line2         string `json:"line_2,omitempty" validate:"omitempty,max=100" example:"Apto 123"`    // Linha 2 do endereço
	District      string `json:"district" validate:"required,max=50" example:"Centro"`                // Bairro
	PostalCode    string `json:"postal_code" validate:"required,cep,max=10" example:"12345678"`       // CEP
	PostalCodeAlt string `json:"postalCode,omitempty" validate:"omitempty,max=10" example:"12345678"` // CEP alternativo
}

type CreateCustomerDTOOutput struct {
	ID string `json:"id" example:"835524dd-d897-45ed-a68e-2e1a8bf522d7"` // Identificador do cliente
}
