package customer_rede_dto

import "encoding/xml"

type AddConsumerInput struct {
	XMLName       xml.Name `xml:"request"`
	CustomerIDExt string   `xml:"customerIdExt"`
	FirstName     string   `xml:"firstName"`
	LastName      string   `xml:"lastName"`
	Address1      string   `xml:"address1"`
	Address2      string   `xml:"address2,omitempty"`
	City          string   `xml:"city"`
	State         string   `xml:"state"`
	Zip           string   `xml:"zip"`
	Country       string   `xml:"country"`
	Phone         string   `xml:"phone"`
	Email         string   `xml:"email"`
	Dob           string   `xml:"dob"`
	Sex           string   `xml:"sex"`
}

type ConsumerResponseOutput struct {
	XMLName      xml.Name `xml:"api-response"`
	ErrorCode    string   `xml:"errorCode"`
	ErrorMessage string   `xml:"errorMessage"`
	Command      string   `xml:"command"`
	Time         string   `xml:"time"`
	Result       struct {
		CustomerID string `xml:"customerId"`
	} `xml:"result"`
}

type ConsumerResponseErrorOutput struct {
	XMLName      xml.Name `xml:"api-response"`
	ErrorCode    string   `xml:"errorCode"`
	ErrorMessage string   `xml:"errorMessage"`
	Time         string   `xml:"time"`
}
