package shortener

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	shortenerAddress  string
	cacheTTLURLCreate time.Duration
	cacheTTLURLVisit  time.Duration
}

func LoadConfig() (Config, error) {
	address := os.Getenv("SHORTENER_ADDRESS")
	if address == "" {
		return Config{}, errors.New("SHORTENER_ADDRESS is empty")
	}

	cacheTTLURLCreateToken := os.Getenv("CACHE_TTL_URL_CREATE")
	if cacheTTLURLCreateToken == "" {
		return Config{}, errors.New("CACHE_TTL_URL_CREATE is empty")
	}
	cacheTTLURLCreate, err := strconv.ParseUint(cacheTTLURLCreateToken, 10, 64)
	if err != nil {
		return Config{}, fmt.Errorf("CACHE_TTL_URL_CREATE has invalid ttl number: %s", cacheTTLURLCreateToken)
	}

	cacheTTLURLVisitToken := os.Getenv("CACHE_TTL_URL_VISIT")
	if cacheTTLURLVisitToken == "" {
		return Config{}, errors.New("CACHE_TTL_URL_VISIT is empty")
	}
	cacheTTLURLVisit, err := strconv.ParseUint(cacheTTLURLVisitToken, 10, 64)
	if err != nil {
		return Config{}, fmt.Errorf("CACHE_TTL_URL_VISIT has invalid ttl number: %s", cacheTTLURLVisitToken)
	}

	conf := Config{
		shortenerAddress:  address,
		cacheTTLURLCreate: time.Duration(cacheTTLURLCreate) * time.Second,
		cacheTTLURLVisit:  time.Duration(cacheTTLURLVisit) * time.Second,
	}
	return conf, nil
}
