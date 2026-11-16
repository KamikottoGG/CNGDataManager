package sheets

import (
	"encoding/json/v2"
	"fmt"
	"strings"
	"time"
)

type BeneficiarySex string

const (
	MALE   BeneficiarySex = "MALE"
	FEMALE BeneficiarySex = "FEMALE"
)

type BeneficiaryPassport struct {
	Number       NullString `json:"number"`
	Issuer       NullString `json:"issuer"`
	IssueDate    DateOnly   `json:"issueDate"`
	IssueCountry NullString `json:"issueCountry"`
}

type NullString string

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns == "" {
		return []byte("null"), nil
	}
	if ns == "-" {
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
	return []byte(fmt.Sprintf("%q", t.Format("2006-01-02"))), nil
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

	for _, layout := range SupportedDateLayouts {
		if parsed, err := time.Parse(layout, s); err == nil {
			*d = DateOnly(parsed)
			return nil
		}
	}

	return fmt.Errorf("incorrect date format %q, expected one of DD/MM/YYYY, YYYY-MM-DD, etc.", s)
}

type RecordMapper[T any] func(map[string]string) (T, error)

type RowFilter func(map[string]string) bool

type RowGetter func() ([]map[string]string, error)

type Cacher[T any] interface {
	Save(records []T) error
}

type Record interface {
	GetIndicatorDate() DateOnly
}
