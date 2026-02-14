package config

import (
	"os"
)

type Config struct {
	Theme string
}

func Load() *Config {
	return &Config{
		Theme: os.Getenv("TFCDASH_THEME"),
	}
}
