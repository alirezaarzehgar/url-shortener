package http

import (
	"errors"
	"os"
)

type Config struct {
	shortenerAddress string
}

func LoadConfig() (Config, error) {
	address := os.Getenv("TRANSPORT_HTTP_RUNNING_ADDRESS")
	if address == "" {
		return Config{}, errors.New("TRANSPORT_HTTP_RUNNING_ADDRESS is empty")
	}
	conf := Config{shortenerAddress: address}
	return conf, nil
}
