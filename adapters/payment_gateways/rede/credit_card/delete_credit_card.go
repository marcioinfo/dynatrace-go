package card_rede_dto

import "encoding/xml"

type DeleteCardDTO struct {
	XMLName    xml.Name `xml:"request"`
	CustomerID string   `xml:"customerId"`
	Token      string   `xml:"token"`
}
