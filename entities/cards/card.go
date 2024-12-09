package cards

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/layer-tools/middlewares"
	"github.com/google/uuid"
)

type Card struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	Holder      string `json:"holder"`
	Brand       string `json:"brand"`
	Fingerprint string `json:"fingerprint,omitempty"`
	FirstDigits string `json:"first_digits"`
	LastDigits  string `json:"last_digits"`
	ExpMonth    string `json:"exp_month"`
	ExpYear     string `json:"exp_year"`
	DeletedAt   datetypes.CustomDateTime
	CreatedAt   datetypes.CustomDateTime
	UpdatedAt   datetypes.CustomDateTime
	ApmLink     middlewares.ApmInfoSend `json:"apmLink"`
}

var CardsAttributes = map[string]string{
	"TableName":    "Cartões",
	"id":           "Id",
	"holder":       "Dono do cartão",
	"brand":        "Marca",
	"first_digits": "Primeiros dígitos",
	"last_disgits": "Últimos dígitos",
	"expire_month": "Mês de vencimento",
	"expire_year":  "Ano de vencimento",
	"created_at":   "Data de criação",
	"updated_at":   "Data da última alteração",
}

func NewCardApplication() *Card {
	return &Card{}
}

func (c *Card) InitID() {
	c.ID = uuid.New().String()
}

func (c *Card) GenerateFingerprint(cardNumber string) string {
	salt := os.Getenv("SALT_FINGERPRINT")
	if salt == "" {
		salt = "secret123"
	}

	data := c.CustomerID + cardNumber + salt
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

func (c *Card) IsValid() error {
	if c.CustomerID == "" {
		return fmt.Errorf("CustomerID é obrigatório")
	}
	if c.Holder == "" {
		return fmt.Errorf("Holder é obrigatório")
	}
	if c.Fingerprint == "" {
		return fmt.Errorf("Fingerprint é obrigatório")
	}
	if c.FirstDigits == "" {
		return fmt.Errorf("FirstDigits é obrigatório")
	}
	if c.LastDigits == "" {
		return fmt.Errorf("LastDigits é obrigatório")
	}
	if c.ExpMonth == "" {
		return fmt.Errorf("ExpMonth é obrigatório")
	}
	if c.ExpYear == "" {
		return fmt.Errorf("ExpYear é obrigatório")
	}
	return nil
}
