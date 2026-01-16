package cmd

import (
	"log/slog"
	"os"

	"github.com/alirezaarzehgar/url-shortener/internal/database"
	"github.com/alirezaarzehgar/url-shortener/internal/database/inmem"
	"github.com/alirezaarzehgar/url-shortener/internal/database/scylladb"
	shortenerService "github.com/alirezaarzehgar/url-shortener/internal/service/shortener"
	transportHttp "github.com/alirezaarzehgar/url-shortener/internal/transport/http"
)

func chooseDatabaseKeyPool(choice database.Choice) database.KeyPool {
	switch choice {
	case database.DatabaseScyllaDB:
		conf, err := scylladb.LoadConfig()
		if err != nil {
			slog.Error("failed to load config", "error", err)
			os.Exit(1)
		}

		connection, err := scylladb.Connect(conf)
		if err != nil {
			slog.Error("failed connecting to database", "error", err)
			os.Exit(1)
		}

		return scylladb.NewKeyPool(connection, conf)
	case database.DatabaseInMemory:
		return inmem.NewKeyPool()
	default:
		panic("invalid")
	}
}

func chooseDatabaseURLShortener(choice database.Choice) database.URLShortener {
	switch choice {
	case database.DatabaseScyllaDB:
		conf, err := scylladb.LoadConfig()
		if err != nil {
			slog.Error("failed to load config", "error", err)
			os.Exit(1)
		}

		connection, err := scylladb.Connect(conf)
		if err != nil {
			slog.Error("failed connecting to database", "error", err)
			os.Exit(1)
		}

		return scylladb.NewShortener(connection)
	case database.DatabaseInMemory:
		return inmem.NewShortener()
	default:
		panic("invalid")
	}
}

func Server(args []string) {
	dbConf, err := database.LoadConfig()
	if err != nil {
		slog.Error("failed to load database config", "error", err)
		os.Exit(1)
	}

	shortenerDB := chooseDatabaseURLShortener(dbConf.ShortenerDBChoice)
	keypoolDB := chooseDatabaseKeyPool(dbConf.KeypoolDBChoice)

	serviceConf, err := shortenerService.LoadConfig()
	if err != nil {
		slog.Error("failed to load database config", "error", err)
		os.Exit(1)
	}

	shortener := shortenerService.New(serviceConf, shortenerDB, keypoolDB)

	transportConf, err := transportHttp.LoadConfig()
	if err != nil {
		slog.Error("failed to load transport config", "error", err)
		os.Exit(1)
	}

	slog.Info("run server")
	t := transportHttp.New(transportConf, shortener)
	if err := t.Run(); err != nil {
		slog.Error("failed to run transport", "error", err)
		os.Exit(1)
	}
}
