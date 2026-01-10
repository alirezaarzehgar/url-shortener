package scylladb

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	cluster []string
}

func LoadConfig() (Config, error) {
	cluster := strings.Fields(os.Getenv("SCYLLADB_CLUSTER"))
	if len(cluster) == 0 {
		return Config{}, errors.New("empty cluster config")
	}
	return Config{cluster: cluster}, nil
}
