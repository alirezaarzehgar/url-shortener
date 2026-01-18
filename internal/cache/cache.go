package cache

import (
	"net/url"
	"time"
)

type Cache interface {
	SetOriginalURL(key string, originalURL url.URL, ttl time.Duration) error
	GetOriginalURL(key string) (url.URL, error)
}
