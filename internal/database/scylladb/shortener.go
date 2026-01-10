package scylladb

import "github.com/alirezaarzehgar/url-shortener/internal/database"

type URLShortenerImpl struct {
	ScyllaDBConnection
}

func (us URLShortenerImpl) Create(key database.URLKey, originalURL database.URL) error {
	panic("not implemented")
}

func (us URLShortenerImpl) Lookup(key database.URLKey) (database.URL, error) {
	panic("not implemented")
}

func NewShortener(conn ScyllaDBConnection) database.URLShortener {
	return URLShortenerImpl{
		ScyllaDBConnection: conn,
	}
}
