package sheets

import (
	"encoding/json"
	"os"
)

type NullString string

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

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(ns))
}
