package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alirezaarzehgar/url-shortener/internal/database/scylladb"
	"github.com/alirezaarzehgar/url-shortener/internal/logger/slogger"
)

func migrateHelp() {
	fmt.Printf("Usage: %s <command>\n", os.Args[0])
	fmt.Println("Commands: ")
	fmt.Println("	scylladb Migrate scylladb database")
	os.Exit(1)
}

func Migratoin(args []string) {
	if len(args) < 2 {
		migrateHelp()
	}

	switch args[1] {
	case "scylladb", "scylla":
		scyllaDBMigration(args[1:])
	default:
		migrateHelp()
	}
}

func scyllaDBMigration(args []string) {
	nLogger := slogger.New()

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	uri := fs.String("uri", "localhost:9042", "scylladb uri")
	migPath := fs.String("migrations", "./internal/database/scylladb/migrations/", "path to migrations directory for migrator")
	clusterSize := fs.Uint("cluster-size", 1, "size of cluster")
	fs.Parse(args[1:])

	migrator, err := scylladb.NewMigrator(*uri, *migPath, *clusterSize)
	if err != nil {
		nLogger.Error("failed to init migrator", "error", err, "uri", *uri, "cluster-size", *clusterSize)
		os.Exit(1)
	}

	if err := migrator.Run(); err != nil {
		nLogger.Error("failed to run migrator", "error", err)
		os.Exit(1)
	}
}
