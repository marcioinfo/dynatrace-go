package customer_pagarme_dto

type AddCustomerInput struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Code         string  `json:"code"`
	Document     string  `json:"document"`
	DocumentType string  `json:"document_type"`
	Type         string  `json:"type"`
	Gender       string  `json:"gender"`
	Address      Address `json:"address"`
	Phones       Phones  `json:"phones"`
	Birthdate    string  `json:"birthdate"`
}

type Address struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
	Line1   string `json:"line_1"`
	Line2   string `json:"line_2"`
}

type Phones struct {
	HomePhone   *Phone `json:"home_phone,omitempty"`
	MobilePhone *Phone `json:"mobile_phone,omitempty"`
}

type Phone struct {
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	Number      string `json:"number"`
}

type AddCustomerOutput struct {
	CustomerToken string  `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Code          string  `json:"code"`
	Document      string  `json:"document"`
	DocumentType  string  `json:"document_type"`
	Type          string  `json:"type"`
	Gender        string  `json:"gender"`
	Delinquent    bool    `json:"delinquent"`
	Address       Address `json:"address"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Birthdate     string  `json:"birthdate"`
	Phones        Phones  `json:"phones"`
	Metadata      struct {
		Id      string `json:"id"`
		Company string `json:"company"`
	} `json:"metadata"`
}
