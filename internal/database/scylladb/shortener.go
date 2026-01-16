package scylladb

import (
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/alirezaarzehgar/url-shortener/internal/database"
	"github.com/gocql/gocql"
)

type URLShortenerImpl struct {
	ScyllaDBConnection
}

func (us URLShortenerImpl) Create(key database.URLKey, originalURL url.URL, ttl uint) error {
	err := us.session.Query(
		`INSERT INTO urlshortener.urls (key, original_url, created_at) VALUES (?, ?, ?) USING TTL ?`,
		key, originalURL.String(), time.Now(), ttl,
	).Exec()
	if err != nil {
		slog.Error("failed to save new short url", "error", err, "url", originalURL.String(), "ttl", ttl)
		return fmt.Errorf("failed to save new short link: %w", err)
	}
	return nil
}

func (us URLShortenerImpl) Lookup(key database.URLKey) (url.URL, error) {
	var originalURL string
	err := us.session.Query(`SELECT original_url FROM urlshortener.urls WHERE key = ?`, key).Scan(&originalURL)

	if err != nil {
		slog.Error("failed to get original url", "error", err, "key", key)
		if err == gocql.ErrNotFound {
			return url.URL{}, database.Err{
				Err:    err,
				Msg:    "original url not found",
				Status: database.NotFoundError,
			}
		} else {
			return url.URL{}, database.Err{
				Err:    err,
				Msg:    "database unable to respond",
				Status: database.InternalError,
			}
		}
	}

	u, err := url.Parse(originalURL)
	if err != nil {
		slog.Error("failed to parse retrieved url from database", "error", err)
		return url.URL{}, database.Err{
			Err:    err,
			Msg:    "database unable to respond",
			Status: database.InternalError,
		}
	}

	return *u, nil
}

func NewShortener(conn ScyllaDBConnection) database.URLShortener {
	return URLShortenerImpl{
		ScyllaDBConnection: conn,
	}
}
