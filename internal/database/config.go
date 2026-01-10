package database

import (
	"fmt"
	"os"
)

type Choice uint

const (
	databaseChoiceStart Choice = iota
	DatabaseScyllaDB
	DatabaseInMemory
	databaseChoiceEnd
)

func LoadChoice(choice string) (Choice, error) {
	switch choice {
	case "scylla", "scylladb":
		return DatabaseScyllaDB, nil
	case "in-memory", "inmemory", "test":
		return DatabaseInMemory, nil
	default:
		return 0, fmt.Errorf("unimplemented database choice: %s", choice)
	}
}

type Config struct {
	ShortenerDBChoice Choice
	KeypoolDBChoice   Choice
}

func LoadConfig() (Config, error) {
	shortenerChoice, err := LoadChoice(os.Getenv("SHORTENER_DATABASE_CHOICE"))
	if err != nil {
		return Config{}, fmt.Errorf("failed to load shortener choice: %w", err)
	}
	keypoolChoice, err := LoadChoice(os.Getenv("KEYPOOL_DATABASE_CHOICE"))
	if err != nil {
		return Config{}, fmt.Errorf("failed to load keypool choice: %w", err)
	}

	conf := Config{
		ShortenerDBChoice: shortenerChoice,
		KeypoolDBChoice:   keypoolChoice,
	}
	return conf, nil
}
