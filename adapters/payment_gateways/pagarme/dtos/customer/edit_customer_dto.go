package customer_pagarme_dto

type EditCustomerDTO struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Phones    Phones `json:"phones"`
	BirthDate string `json:"birthdate,omitempty"`
	Document  string `json:"document,omitempty"`
	Type      string `json:"type,omitempty"`
}
