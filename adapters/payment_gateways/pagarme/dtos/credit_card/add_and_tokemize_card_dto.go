package card_pagarme_dto

type AddAndTokemizeCardInput struct {
	Number          string         `json:"number"`
	HolderName      string         `json:"holder_name"`
	HolderDocument  string         `json:"holder_document,omitempty"`
	ExpirationMonth string         `json:"exp_month"`
	ExpirationYear  string         `json:"exp_year"`
	CVV             string         `json:"cvv"`
	Label           string         `json:"label"`
	BillingAddress  BillingAddress `json:"billing_address"`
	Options         Options        `json:"options"`
}

type Options struct {
	VerifyCard bool `json:"verify_card"`
}

type BillingAddress struct {
	Line1   string `json:"line_1"`
	Line2   string `json:"line_2"`
	ZipCode string `json:"zip_code"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type AddAndTokemizeCardOutput struct {
	CardToken  string `json:"id"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
	LastDigits string `json:"last_four_digits"`
	HolderName string `json:"holder_name"`
}
