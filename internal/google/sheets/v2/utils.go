package sheets

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

func text(msg ...any) string {
	return fmt.Sprintf("%+v", msg...)
}

func ExtractValuesFromRows(rows *sheets.ValueRange, offset int) []map[string]string {

	var values = make([]map[string]string, 0, len(rows.Values))

	columnHeaders := rows.Values[0]

	for _, row := range rows.Values[offset:] {

		mappedObj := make(map[string]string)

		for i, value := range row {
			mappedObj[text(columnHeaders[i])] = text(value)
		}
		values = append(values, mappedObj)
	}

	return values
}

var SupportedDateLayouts = []string{
	"02/01/2006",
	"2/1/2006",
	"02-01-2006",
	"2-1-2006",
	"02.01.2006",
	"2.1.2006",
	"2.1-2006",
	"2006/01/02",
	"2006-01-02",
	"01-02-2006",
	"1-2-2006",
	"01/02/2006",
	"02.01.2006 15:04:05",
}

func ParseSheetDate(s string) (DateOnly, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return DateOnly{}, nil
	}
	for _, layout := range SupportedDateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return DateOnly(t), nil
		}
	}
	return DateOnly{}, fmt.Errorf("incorrect date %q, expected DD/MM/YYYY", s)
}

func ParseYesNo(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "так", "yes", "true", "+", "1":
		return true
	}
	return false
}

func ParseSex(s string) (BeneficiarySex, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "чоловіча", "чоловік", "ч", "м", "male", "man", "х":
		return MALE, nil
	case "жіноча", "жінка", "ж", "f", "female", "д":
		return FEMALE, nil
	}
	return "", fmt.Errorf("unknown sex value %q", s)
}
