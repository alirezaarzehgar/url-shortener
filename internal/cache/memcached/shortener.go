package memcached

import (
	"net/url"
	"time"

	"github.com/alirezaarzehgar/url-shortener/internal/cache"
)

type Memcached struct {
}

func (c Memcached) SetOriginalURL(key string, originalURL url.URL, ttl time.Duration) error {
	panic("not implemented")
}

func (c Memcached) GetOriginalURL(key string) (url.URL, error) {
	panic("not implemented")
}

func New() cache.Cache {
	return Memcached{}
}
