package sheets

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONCacher[T any] struct {
	FilePath string
	mu       sync.RWMutex
}

func NewJSONCacher[T any](fp string) *JSONCacher[T] {
	return &JSONCacher[T]{FilePath: fp}
}

func (c *JSONCacher[T]) Save(obj []T) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	f, err := os.Create(c.FilePath + ".tmp")
	if err != nil {

		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(obj)
	if err != nil {
		f.Close()
		os.Remove(c.FilePath + ".tmp")
		return err
	}
	f.Close()
	return os.Rename(c.FilePath+".tmp", c.FilePath)
}
