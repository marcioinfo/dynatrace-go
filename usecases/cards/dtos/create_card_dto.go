package dtos

import "github.com/adhfoundation/layer-tools/datetypes"

type CreateCardDTOInput struct {
	CustomerID     string          `json:"customer_id" validate:"omitempty,uuid4" example:"835524dd-d897-45ed-a68e-2e1a8bf52314"` // Identificador do cliente
	Holder         string          `json:"holder" validate:"required,max=255" example:"Joao Silva"`                               // Nome do portador do cartão
	Brand          string          `json:"brand,omitempty" validate:"omitempty,max=50" example:"visa"`                            // Bandeira do cartão
	Number         string          `json:"number" validate:"required,numeric,max=20" example:"4111111111111111"`                  // Número do cartão
	CVV            string          `json:"cvv" validate:"required,numeric,max=4" example:"123"`                                   // Código de segurança do cartão
	ExpMonth       string          `json:"exp_month" validate:"required,numeric,month,max=2" example:"12"`                        // Mês de expiração do cartão
	ExpYear        string          `json:"exp_year" validate:"required,numeric,len=2" example:"25"`                               // Ano de expiração do cartão
	BillingAddress *BillingAddress `json:"billing_address" validate:"required"`                                                   // Endereço de cobrança do cartão
	Customer       *Customer       `json:"customer,omitempty"`                                                                    // Cliente do cartão
}

type Customer struct {
	FirstName string               `json:"first_name" validate:"required,max=50" example:"João"`                                // Nome do cliente
	LastName  string               `json:"last_name" validate:"required,max=50" example:"Silva"`                                // Sobrenome do cliente
	Document  string               `json:"document" validate:"required,cpf_or_cnpj,max=14" example:"12345678901"`               // CPF ou CNPJ
	BirthDate datetypes.CustomDate `json:"birth_date" validate:"required" example:"2000-11-06"`                                 // Data de nascimento
	Gender    string               `json:"gender" validate:"required,oneof=male female" example:"male || female"`               // Gênero
	Phone     string               `json:"phone" validate:"required,numeric,max=13" example:"11999999999"`                      // Telefone
	Email     string               `json:"email" validate:"required,email,max=100" example:"teste@gmail.com"`                   // E-mail
	Address   BillingAddress       `json:"address" validate:"required"`                                                         // Endereço de cobrança
	ServiceID string               `json:"service_id" validate:"required,uuid4" example:"835524dd-d897-45ed-a68e-2e1a8bf52314"` // ID do serviço
}

type BillingAddress struct {
	State         string `json:"state" validate:"required,states,max=2" example:"SP"`                 // Estado
	City          string `json:"city" validate:"required,max=50" example:"São Paulo"`                 // Cidade
	Country       string `json:"country" validate:"required,iso3166_1_alpha2,len=2" example:"BR"`     // País
	Line1         string `json:"line_1" validate:"required,max=100" example:"Rua João Silva, 123"`    // Linha 1 do endereço
	Line2         string `json:"line_2,omitempty" validate:"omitempty,max=100" example:"Apto 123"`    // Linha 2 do endereço
	District      string `json:"district,omitempty" validate:"omitempty,max=50" example:"Centro"`     // Bairro
	PostalCode    string `json:"postal_code" validate:"required,cep,max=10" example:"12345678"`       // CEP
	PostalCodeAlt string `json:"postalCode,omitempty" validate:"omitempty,max=10" example:"12345678"` // CEP alternativo
}

type CreateCardDTOOutput struct {
	ID         string `json:"id" validate:"omitempty,uuid4" example:"835524dd-d897-45ed-a68e-2e1a8bf522d7"`          // Identificador do cartão
	CustomerID string `json:"customer_id" validate:"omitempty,uuid4" example:"835524dd-d897-45ed-a68e-2e1a8bf52314"` // Identificador do cliente
}
