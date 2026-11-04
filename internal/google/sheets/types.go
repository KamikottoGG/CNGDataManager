package sheets

type BeneficiarySex string

const (
	MALE   BeneficiarySex = "MAN"
	FEMALE BeneficiarySex = "FEMALE"
)

type BeneficiaryPassport struct {
	Number       string   `json:"number"`
	Issuer       string   `json:"issuer"`
	IssueDate    DateOnly `json:"issueDate"`
	IssueCountry string   `json:"issueCountry"`
}
