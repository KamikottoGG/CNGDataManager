package sheets

import (
	"encoding/json/v2"
	"fmt"
	"strings"
	"time"
)

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

type NullString string

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(ns))
}

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%q", t.Format(DATE_LAYOUT))), nil
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))

	if s == "null" || s == `""` || s == "" {
		*d = DateOnly(time.Time{})
		return nil
	}

	s = strings.Trim(s, `"`)
	s = strings.TrimSpace(s)

	if s == "" || s == "-" {
		*d = DateOnly(time.Time{})
		return nil
	}

	parsed, err := time.Parse(DATE_LAYOUT, s)
	if err != nil {
		return fmt.Errorf("incorrect date format %q, expected %s: %w", s, DATE_LAYOUT, err)
	}

	*d = DateOnly(parsed)
	return nil
}
