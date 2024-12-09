package card_rede_dto

import "encoding/xml"

type AddAndTokemizeCardInput struct {
	XMLName               xml.Name `xml:"request"`
	ExpirationMonth       string   `xml:"expirationMonth"`
	CustomerID            string   `xml:"customerId"`
	CreditCardNumber      string   `xml:"creditCardNumber"`
	ExpirationYear        string   `xml:"expirationYear"`
	BillingName           string   `xml:"billingName"`
	BillingAddress1       string   `xml:"billingAddress1"`
	BillingAddress2       string   `xml:"billingAddress2"`
	BillingCity           string   `xml:"billingCity"`
	BillingState          string   `xml:"billingState"`
	BillingZip            string   `xml:"billingZip"`
	BillingCountry        string   `xml:"billingCountry"`
	BillingPhone          string   `xml:"billingPhone"`
	BillingEmail          string   `xml:"billingEmail"`
	OnFileEndDate         string   `xml:"onFileEndDate,omitempty"`
	OnFilePermissions     string   `xml:"onFilePermissions,omitempty"`
	OnFileComment         string   `xml:"onFileComment,omitempty"`
	OnFileMaxChargeAmount string   `xml:"onFileMaxChargeAmount,omitempty"`
}

type ZeroDollarCard struct {
	TransactionDetail TransactionDetail `xml:"transactionDetail"`
	ProcessorID       string            `xml:"processorID"`
	ReferenceNum      string            `xml:"referenceNum"`
	SaveOnFile        SaveOnFile        `xml:"saveOnFile"`
}

type SaveOnFile struct {
	CustomerToken string `xml:"customerToken"`
}
type Order struct {
	ZeroDollarCard ZeroDollarCard `xml:"zeroDollar"`
	XMLName        xml.Name       `xml:"order"`
}

type TransactionDetail struct {
	PayType PayType `xml:"payType"`
}

type PayType struct {
	CreditCard CreditCard `xml:"creditCard"`
}

type CreditCard struct {
	Number    string `xml:"number"`
	ExpMonth  string `xml:"expMonth"`
	ExpYear   string `xml:"expYear"`
	CvvNumber string `xml:"cvvNumber"`
}

type AddAndTokemizeCardSuccessOutput struct {
	XMLName      xml.Name `xml:"api-response"`
	ErrorCode    string   `xml:"errorCode"`
	ErrorMessage string   `xml:"errorMessage"`
	Command      string   `xml:"command"`
	Time         string   `xml:"time"`
	Result       struct {
		Token string `xml:"token"`
	} `xml:"result"`
}

type AddAndTokemizeCardWithZeroDollar struct {
	XMLName                  xml.Name `xml:"transaction-response"`
	AuthCode                 string   `xml:"authCode"`
	OrderID                  string   `xml:"orderID"`
	ReferenceNum             string   `xml:"referenceNum"`
	TransactionID            string   `xml:"transactionID"`
	TransactionTimestamp     string   `xml:"transactionTimestamp"`
	ResponseCode             string   `xml:"responseCode"`
	ResponseMessage          string   `xml:"responseMessage"`
	AvsResponseCode          string   `xml:"avsResponseCode"`
	CvvResponseCode          string   `xml:"cvvResponseCode"`
	ProcessorCode            string   `xml:"processorCode"`
	ProcessorMessage         string   `xml:"processorMessage"`
	ProcessorName            string   `xml:"processorName"`
	ErrorMessage             string   `xml:"errorMessage"`
	ProcessorTransactionID   string   `xml:"processorTransactionID"`
	ProcessorReferenceNumber string   `xml:"processorReferenceNumber"`
	CreditCardScheme         string   `xml:"creditCardScheme"`
	SaveFile                 SaveFile `xml:"save-on-file"`
	ErrorCode                string   `xml:"errorCode"`
}

type SaveFile struct {
	Token string `xml:"token"`
}

type AddAndTokemizeCardResponseErrorOutput struct {
	XMLName      xml.Name `xml:"api-response"`
	ErrorCode    string   `xml:"errorCode"`
	ErrorMessage string   `xml:"errorMessage"`
	Time         string   `xml:"time"`
}
