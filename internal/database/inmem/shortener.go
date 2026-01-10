package inmem

import (
	"sync"

	"github.com/alirezaarzehgar/url-shortener/internal/database"
)

type URLShortenerImpl struct {
	mutex   *sync.RWMutex
	storage map[database.URLKey]database.URL
}

func (us URLShortenerImpl) Create(key database.URLKey, originalURL database.URL) error {
	us.mutex.Lock()
	us.storage[key] = originalURL
	us.mutex.Unlock()
	return nil
}

func (us URLShortenerImpl) Lookup(key database.URLKey) (database.URL, error) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	originalURL, ok := us.storage[key]
	if !ok {
		return nil, database.Err{Msg: "not found", Status: database.NotFoundError}
	}
	return originalURL, nil
}

func NewShortener() database.URLShortener {
	return URLShortenerImpl{
		mutex:   &sync.RWMutex{},
		storage: make(map[database.URLKey]database.URL),
	}
}
