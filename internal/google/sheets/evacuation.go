package sheets

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DATE_LAYOUT = "2006/01/02"
)

type DateOnly time.Time

type Evacuation struct {
	Date DateOnly `json:"date"`

	Name       string     `json:"name"`
	Surname    string     `json:"surname"`
	Patronymic NullString `json:"patronymic"`

	BirthDate      DateOnly            `json:"birthDate"`
	Age            int                 `json:"age"`
	Sex            BeneficiarySex      `json:"sex"`
	Disability     bool                `json:"disability"`
	Passport       BeneficiaryPassport `json:"passport"`
	TaxId          NullString          `json:"taxId"`
	Region         NullString          `json:"region"`
	District       NullString          `json:"district"`
	Community      NullString          `json:"community"`
	SettlementType NullString          `json:"settlementType"`
	Settlement     NullString          `json:"settlement"`

	EvacuatedTo NullString `json:"evacuatedTo"`
	PhoneNumber NullString `json:"phoneNumber"`
	Comment     NullString `json:"comment"`
}

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

const SHEET_DATE_LAYOUT = "02/01/2006"

func parseSheetDate(s string) (DateOnly, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return DateOnly{}, nil
	}
	for _, layout := range []string{
		SHEET_DATE_LAYOUT,
		"2/1/2006",
		"02-01-2006",
		"2-1-2006",
		"02.01.2006",
		"2.1.2006",
		DATE_LAYOUT,
		"2006-01-02",

		"01-02-2006",
		"1-2-2006",
		"01/02/2006",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return DateOnly(t), nil
		}
	}
	return DateOnly{}, fmt.Errorf("incorrect date %q, expected DD/MM/YYYY", s)
}

func parseYesNo(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "так", "yes", "true", "+", "1":
		return true
	}
	return false
}

func parseSex(s string) (BeneficiarySex, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "чоловіча", "чоловік", "ч", "м", "male", "man":
		return MALE, nil
	case "жіноча", "жінка", "ж", "f", "female":
		return FEMALE, nil
	}
	return "", fmt.Errorf("unknown sex value %q", s)
}

func (e *Evacuation) MapObject(object map[string]string) error {
	get := func(key string) string { return strings.TrimSpace(object[key]) }

	var err error

	if e.Date, err = parseSheetDate(get("Дата [ДД/ММ/РРРР]")); err != nil {
		return fmt.Errorf("evacuation date: %w", err)
	}

	e.Surname = get("Прізвище")
	e.Name = get("Ім'я")
	e.Patronymic = NullString(get("По батькові"))

	if e.BirthDate, err = parseSheetDate(get("Дата народження [ДД/ММ/РРРР]")); err != nil {
		return fmt.Errorf("birth date: %w", err)
	}

	if age, convErr := strconv.Atoi(get("Вік")); convErr == nil {
		e.Age = age
	} else if bt := time.Time(e.BirthDate); !bt.IsZero() {
		now := time.Now()
		e.Age = now.Year() - bt.Year()
		if now.YearDay() < bt.YearDay() {
			e.Age--
		}
	}

	if e.Sex, err = parseSex(get("Стать")); err != nil {
		return fmt.Errorf("sex: %w", err)
	}
	e.Disability = parseYesNo(get("Інвалідність (так/ні)"))

	e.Passport.Number = get("Паспорт")
	e.Passport.Issuer = get("Орган що видав")
	if e.Passport.IssueDate, err = parseSheetDate(get("Дата видачі")); err != nil {
		return fmt.Errorf("passport issue date: %w", err)
	}
	e.Passport.IssueCountry = get("Країна видачі")

	e.TaxId = NullString(get("ІПН"))
	e.Region = NullString(get("Звідки область"))
	e.District = NullString(get("Звідки район"))
	e.Community = NullString(get("Звідки Громада"))
	e.SettlementType = NullString(get("Тип населеного пункту"))
	e.Settlement = NullString(get("Звідки Наcпункт"))

	e.EvacuatedTo = NullString(get("Куди евакуювалися"))
	e.PhoneNumber = NullString(get("Номер телефону"))
	e.Comment = NullString(get("Коментар"))

	return nil
}
