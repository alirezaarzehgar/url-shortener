package scylladb

import (
	"fmt"

	"github.com/gocql/gocql"
)

type ScyllaDBConnection struct {
	session *gocql.Session
}

func Connect(conf Config) (ScyllaDBConnection, error) {
	cluster := gocql.NewCluster(conf.cluster...)
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return ScyllaDBConnection{}, fmt.Errorf("failed to create scylladb session: %w", err)
	}

	return ScyllaDBConnection{
		session: session,
	}, nil
}
