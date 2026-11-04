package sheets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

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
