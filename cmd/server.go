package cmd

import (
	"fmt"
	"os"

	"github.com/alirezaarzehgar/url-shortener/internal/database"
	"github.com/alirezaarzehgar/url-shortener/internal/database/inmem"
	"github.com/alirezaarzehgar/url-shortener/internal/database/scylladb"
	"github.com/alirezaarzehgar/url-shortener/internal/logger"
	"github.com/alirezaarzehgar/url-shortener/internal/logger/slogger"
	shortenerService "github.com/alirezaarzehgar/url-shortener/internal/service/shortener"
	transportHttp "github.com/alirezaarzehgar/url-shortener/internal/transport/http"
)

func chooseDatabaseKeyPool(choice database.Choice) (database.KeyPool, error) {
	switch choice {
	case database.DatabaseScyllaDB:
		conf, err := scylladb.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}

		connection, err := scylladb.Connect(conf)
		if err != nil {
			return nil, fmt.Errorf("failed connecting to database: %w", err)
		}

		return scylladb.NewKeyPool(connection, conf), nil
	case database.DatabaseInMemory:
		return inmem.NewKeyPool(), nil
	default:
		return nil, fmt.Errorf("invalid database type: %v", choice)
	}
}

func chooseDatabaseURLShortener(choice database.Choice, log logger.Logger) (database.URLShortener, error) {
	switch choice {
	case database.DatabaseScyllaDB:
		conf, err := scylladb.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}

		connection, err := scylladb.Connect(conf)
		if err != nil {
			return nil, fmt.Errorf("failed connecting to database: %w", err)
		}

		return scylladb.NewShortener(connection, log), nil
	case database.DatabaseInMemory:
		return inmem.NewShortener(), nil
	default:
		return nil, fmt.Errorf("invalid database type: %v", choice)
	}
}

func Server(args []string) {
	nLogger := slogger.New()

	dbConf, err := database.LoadConfig()
	if err != nil {
		nLogger.Error("failed to load database config", "error", err)
		os.Exit(1)
	}

	shortenerDB, err := chooseDatabaseURLShortener(dbConf.ShortenerDBChoice, nLogger)
	if err != nil {
		nLogger.Error("failed to choose url shortener database", "error", err)
		os.Exit(1)
	}
	keypoolDB, err := chooseDatabaseKeyPool(dbConf.KeypoolDBChoice)
	if err != nil {
		nLogger.Error("failed to choose key pool database", "error", err)
		os.Exit(1)
	}

	serviceConf, err := shortenerService.LoadConfig()
	if err != nil {
		nLogger.Error("failed to load database config", "error", err)
		os.Exit(1)
	}

	shortener := shortenerService.New(serviceConf, nLogger, shortenerDB, keypoolDB)

	transportConf, err := transportHttp.LoadConfig()
	if err != nil {
		nLogger.Error("failed to load transport config", "error", err)
		os.Exit(1)
	}

	nLogger.Info("run server")
	t := transportHttp.New(transportConf, shortener)
	if err := t.Run(); err != nil {
		nLogger.Error("failed to run transport", "error", err)
		os.Exit(1)
	}
}
