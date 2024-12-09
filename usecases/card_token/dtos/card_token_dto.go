package dtos

type CardTokenOutputDTO struct {
	Gateway string `json:"gateway"`
	Token   string `json:"token"`
}

type CardTokenInputDTO struct {
	CardID          string `json:"card_id"`
	CustomerID      string `json:"customer_id"`
	CustomerToken   string `json:"customer_token"`
	CustomerEmail   string `json:"customer_email" validate:"required,email"`
	CustomerAddress *HolderAddress
	CustomerPhone   string `json:"customer_phone" validate:"required,phone,numeric"`
	Number          string `json:"number" validate:"required,numeric"`
	Holder          string `json:"holder"`
	HolderDocument  string `json:"holder_document"`
	ExpMonth        string `json:"exp_month" validate:"required,month,numeric"`
	ExpYear         string `json:"exp_year" validate:"required,year,numeric"`
	CVV             string `json:"cvv" validate:"required,numeric"`
	Brand           string `json:"brand"`
}

type HolderAddress struct {
	State      string `json:"state" `
	City       string `json:"city" `
	Country    string `json:"country" `
	Line1      string `json:"line_1" `
	Line2      string `json:"line_2" `
	District   string `json:"district" `
	PostalCode string `json:"postal_code" `
}
