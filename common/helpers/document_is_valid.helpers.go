package helpers

import (
	"github.com/klassmann/cpfcnpj"
)

func CNPJIsValid(document string) bool {
	cnpj := cpfcnpj.NewCNPJ(document)
	return cnpj.IsValid()
}

func CPFIsValid(document string) bool {
	cpf := cpfcnpj.NewCPF(document)
	return cpf.IsValid()
}

func CleanDocument(document string) string {
	return cpfcnpj.Clean(document)
}

func DocumentIsValid(document string) bool {
	documentType := CheckDocumentType(document)
	if documentType == "CPF" {
		return CPFIsValid(document)
	} else if documentType == "CNPJ" {
		return CNPJIsValid(document)
	}
	return false
}

func CheckDocumentType(document string) string {
	if len(document) == 11 {
		return "CPF"
	} else if len(document) == 14 {
		return "CNPJ"
	}
	return ""
}
