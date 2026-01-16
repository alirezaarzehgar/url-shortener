package scylladb

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultKeyRange = 50
)

type Config struct {
	cluster   []string
	clusterID int
	keyRange  uint64
}

func LoadConfig() (Config, error) {
	conf := Config{}

	conf.cluster = strings.Fields(os.Getenv("SCYLLADB_CLUSTER"))
	if len(conf.cluster) == 0 {
		return Config{}, errors.New("SCYLLADB_CLUSTER: empty cluster config")
	}

	keyRangeToken := os.Getenv("SCYLLADB_KEY_RANGE")
	if keyRange, _ := strconv.ParseUint(keyRangeToken, 10, 64); keyRange > 0 {
		conf.keyRange = keyRange
	} else {
		conf.keyRange = DefaultKeyRange
	}

	clusterIDToken := os.Getenv("SCYLLADB_CLUSTER_ID")
	if clusterID, err := strconv.Atoi(clusterIDToken); err != nil {
		return Config{}, fmt.Errorf("invaid SCYLLADB_CLUSTER_ID: %w", err)
	} else {
		conf.clusterID = clusterID
	}

	return conf, nil
}
