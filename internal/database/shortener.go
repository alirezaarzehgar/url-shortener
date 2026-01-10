package database

import "net/url"

type URLKey string
type URL *url.URL

type URLShortener interface {
	Create(key URLKey, originalURL URL) error
	Lookup(key URLKey) (URL, error)
}
