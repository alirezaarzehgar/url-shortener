package database

import "net/url"

type URLKey string

type URLShortener interface {
	Create(key URLKey, originalURL url.URL, ttl uint) error
	Lookup(key URLKey) (url.URL, error)
}
