package config

import (
	"os"
	"strings"
)

type Config interface {
	Get(k string) string
}

type EnvConfig struct {
	data map[string]string
}

func NewEnvConfig() *EnvConfig {
	c := new(EnvConfig)

	c.data = make(map[string]string)

	for _, v := range os.Environ() {
		key, value, _ := strings.Cut(v, "=")

		c.data[key] = value
	}
	return c
}

func (c *EnvConfig) Get(k string) string {

	return c.data[k]
}
