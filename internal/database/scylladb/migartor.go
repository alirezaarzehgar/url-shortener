package scylladb

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
)

type Migrator struct {
	ScyllaDBConnection
	clusterSize   uint
	migrationPath string
}

func (m Migrator) Run() error {
	if err := m.runMigrations(); err != nil {
		return fmt.Errorf("failed to apply all migrations: %w", err)
	}
	if err := m.seedKeyPool(); err != nil {
		return fmt.Errorf("failed to seed key pool based on cluster size: %w", err)
	}
	return nil
}

func (m Migrator) runSingleMigration(path string) error {
	migQueries, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	for query := range strings.SplitSeq(string(migQueries), "\n") {
		if len(query) < 2 {
			continue
		}
		err := m.session.Query(query).Exec()
		if err != nil {
			return fmt.Errorf("failed to run query: %w", err)
		}
	}

	return nil
}

func (m Migrator) runMigrations() error {
	de, err := os.ReadDir(m.migrationPath)
	if err != nil {
		return fmt.Errorf("failed to open migartions directory: %s: %w", m.migrationPath, err)
	}

	migFiles := []string{}
	for _, e := range de {
		if _, err := strconv.Atoi(e.Name()[:4]); err != nil {
			continue
		}

		migFiles = append(migFiles, e.Name())
	}
	sort.Slice(migFiles, func(i, j int) bool {
		iN, _ := strconv.Atoi(migFiles[i][:4])
		jN, _ := strconv.Atoi(migFiles[j][:4])
		return iN < jN
	})

	for _, mig := range migFiles {
		migPath := filepath.Join(m.migrationPath, mig)
		if err := m.runSingleMigration(migPath); err != nil {
			return fmt.Errorf("failed to apply migration: %s: %w", mig, err)
		}
	}

	return nil
}

func (m Migrator) seedKeyPool() error {
	b := m.session.NewBatch(gocql.LoggedBatch)
	for id := range m.clusterSize {
		startCounter := id * ((1 << 31) - 1)
		b.Query(`INSERT INTO urlshortener.keypool (key, counter) VALUES (?, ?);`, id, startCounter)
	}
	if err := m.session.ExecuteBatch(b); err != nil {
		return fmt.Errorf("failed to seed key pool: %w", err)
	}

	return nil
}

func NewMigrator(uri string, migrationPath string, clusterSize uint) (Migrator, error) {
	conf := Config{cluster: []string{uri}}
	conn, err := Connect(conf)
	if err != nil {
		return Migrator{}, err
	}

	return Migrator{
		ScyllaDBConnection: conn,
		migrationPath:      migrationPath,
		clusterSize:        clusterSize,
	}, nil
}
