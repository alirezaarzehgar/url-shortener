package local

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/alirezaarzehgar/url-shortener/internal/cache"
)

type LocalCache struct {
	storage *sync.Map
}

type TTLValue struct {
	expiry time.Time
	value  url.URL
}

func (c LocalCache) SetOriginalURL(key string, originalURL url.URL, ttl time.Duration) error {
	ttlv := TTLValue{
		expiry: time.Now().Add(ttl),
		value:  originalURL,
	}
	c.storage.Store(key, ttlv)
	return nil
}

func (c LocalCache) GetOriginalURL(key string) (url.URL, error) {
	v, ok := c.storage.Load(key)
	if !ok {
		err := fmt.Errorf("key not found: %s", key)
		return url.URL{}, cache.Err{Msg: err.Error(), Status: cache.NotFoundError, Err: err}
	}
	ttlv, ok := v.(TTLValue)
	if !ok {
		err := errors.New("internal server error")
		return url.URL{}, cache.Err{Msg: err.Error(), Status: cache.InternalError, Err: err}
	}
	if time.Now().After(ttlv.expiry) {
		c.storage.Delete(key)
		err := fmt.Errorf("key not found: %s", key)
		return url.URL{}, cache.Err{Msg: err.Error(), Status: cache.NotFoundError, Err: err}
	}
	return ttlv.value, nil
}

func New() cache.Cache {
	return LocalCache{
		storage: &sync.Map{},
	}
}
