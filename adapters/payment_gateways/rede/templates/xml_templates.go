package templates

import (
	"fmt"
	"os"
)

func GetTransactionRequestTemplateValue() string {
	template := `<transaction-request>
<version>3.1.1.15</version>
    <verification>
        <merchantId>%s</merchantId>
        <merchantKey>%s</merchantKey>
    </verification>
    {{.Request}}
</transaction-request>`
	return fmt.Sprintf(template, os.Getenv("API_REDE_MERCHANT_ID"), os.Getenv("API_REDE_MERCHANT_KEY"))
}

func GetApiRequestTemplateValue() string {
	template := `<api-request>
<verification>
	<merchantId>%s</merchantId>
	<merchantKey>%s</merchantKey>
</verification>
<command>{{.Command}}</command>
{{.Request}}
</api-request>`
	return fmt.Sprintf(template, os.Getenv("API_REDE_MERCHANT_ID"), os.Getenv("API_REDE_MERCHANT_KEY"))
}
