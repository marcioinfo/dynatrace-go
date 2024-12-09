package customer_rede_dto

import "encoding/xml"

type UpdateCustomerDTO struct {
	XMLName    xml.Name `xml:"request"`
	CustomerID string   `xml:"customerId"`
	FirstName  string   `xml:"firstName"`
	LastName   string   `xml:"lastName"`
	Phone      string   `xml:"phone"`
	Email      string   `xml:"email"`
	Dob        string   `xml:"dob"`
	Sex        string   `xml:"sex"`
}
