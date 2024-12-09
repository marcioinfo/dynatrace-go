package customers

import (
	"fmt"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/layer-tools/middlewares"
	"github.com/google/uuid"
)

type Customer struct {
	ID        string                   `json:"id"`
	Name      string                   `json:"name"`
	Document  string                   `json:"document"`
	BirthDate datetypes.CustomDate     `json:"birth_date"`
	Email     string                   `json:"email"`
	Phone     string                   `json:"phone"`
	Gender    string                   `json:"gender"`
	DeletedAt datetypes.CustomDateTime `json:"deleted_at,omitempty"`
	CreatedAt datetypes.CustomDateTime `json:"created_at"`
	UpdatedAt datetypes.CustomDateTime `json:"updated_at"`
	ApmLink   middlewares.ApmInfoSend  `json:"apmLink"`
}

var CustomerAttributes = map[string]string{
	"TableName":  "Clientes",
	"id":         "Id",
	"name":       "Nome Completo",
	"document":   "Documento",
	"birth_date": "Data de nascimento",
	"email":      "E-mail",
	"phone":      "Telefone",
	"gender":     "Gênero",
	"created_at": "Data de criação",
	"updated_at": "Data da última alteração",
}

func NewCustomer() *Customer {
	return &Customer{}
}

func (c *Customer) InitID() {
	c.ID = uuid.New().String()
}

func (c *Customer) IsValid() error {
	if c.Name == "" {
		return fmt.Errorf("name é obrigatório")
	}
	if c.Document == "" {
		return fmt.Errorf("document é obrigatório")
	}
	if c.BirthDate == (datetypes.CustomDate{}) {
		return fmt.Errorf("birthDate é obrigatório")
	}
	if c.Email == "" {
		return fmt.Errorf("email é obrigatório")
	}
	if c.Phone == "" {
		return fmt.Errorf("phone é obrigatório")
	}

	return nil
}
