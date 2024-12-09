package validators

var STATES = []string{
	"AC",
	"AL",
	"AP",
	"AM",
	"BA",
	"CE",
	"DF",
	"ES",
	"GO",
	"MA",
	"MS",
	"MT",
	"MG",
	"PA",
	"PB",
	"PR",
	"PE",
	"PI",
	"RJ",
	"RN",
	"RS",
	"RO",
	"RR",
	"SC",
	"SP",
	"SE",
	"TO",
}

func IsState(str string) bool {
	for _, v := range STATES {
		if v == str {
			return true
		}
	}

	return false
}
