package shortener

import (
	"errors"
	"os"
)

type Config struct {
	shortenerAddress string
}

func LoadConfig() (Config, error) {
	address := os.Getenv("SHORTENER_ADDRESS")
	if address == "" {
		return Config{}, errors.New("SHORTENER_ADDRESS is empty")
	}
	conf := Config{shortenerAddress: address}
	return conf, nil
}
