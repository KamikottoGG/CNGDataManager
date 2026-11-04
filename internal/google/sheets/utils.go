package sheets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func WriteToJsonFile(filePath string, data any) error {
	f, err := os.Create(filePath + ".tmp")
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)

	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return os.Rename(filePath+".tmp", filePath)
}

func ParseSheetDate(s string) (DateOnly, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return DateOnly{}, nil
	}
	for _, layout := range []string{
		"02/01/2006",
		"2/1/2006",
		"02-01-2006",
		"2-1-2006",
		"02.01.2006",
		"2.1.2006",
		"2006/01/02",
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

func ParseYesNo(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "так", "yes", "true", "+", "1":
		return true
	}
	return false
}

func ParseSex(s string) (BeneficiarySex, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "чоловіча", "чоловік", "ч", "м", "male", "man":
		return MALE, nil
	case "жіноча", "жінка", "ж", "f", "female":
		return FEMALE, nil
	}
	return "", fmt.Errorf("unknown sex value %q", s)
}
