package shortener

import (
	"fmt"
	"net/url"

	"github.com/alirezaarzehgar/url-shortener/internal/database"
	"github.com/alirezaarzehgar/url-shortener/internal/logger"
	"github.com/alirezaarzehgar/url-shortener/internal/service"
)

type URLShortener struct {
	urlResultAddress string
	shortenerDB      database.URLShortener
	keyPoolDB        database.KeyPool
	keyPoolRange     *database.KeyRange
	log              logger.Logger
}

func (us URLShortener) nextKey() (database.URLKey, error) {
	if us.keyPoolRange.Empty() {
		moreKeys, err := us.keyPoolDB.ReserveRange()
		if err != nil {
			return "", fmt.Errorf("failed to reserve new key range from pool: %w", err)
		}
		*us.keyPoolRange = append(*us.keyPoolRange, moreKeys...)
	}
	key, _ := us.keyPoolRange.Pop()
	return database.URLKey(key), nil
}

func (us URLShortener) generateShortURL(key database.URLKey) service.URLKey {
	shortenedURL, _ := url.JoinPath(us.urlResultAddress, string(key))
	return service.URLKey(shortenedURL)
}

func (us URLShortener) CreateShortURL(originalURL url.URL, ttl uint) (service.URLKey, error) {
	key, err := us.nextKey()
	if err != nil {
		return "", fmt.Errorf("failed to get new key: %w", err)
	}

	err = us.shortenerDB.Create(key, originalURL, ttl)
	if err != nil {
		return "", fmt.Errorf("failed store shortened URL: %w", err)
	}

	return us.generateShortURL(key), nil
}

func (us URLShortener) GetOriginalURL(shortURL service.URLKey) (url.URL, error) {
	key := database.URLKey(shortURL)
	originalURL, err := us.shortenerDB.Lookup(key)
	if err != nil {
		if dbErr := err.(database.Err); dbErr.NotFound() {
			return url.URL{}, service.Err{Msg: "short url key not found", Status: service.NotFoundError, Err: dbErr}
		} else {
			return url.URL{}, service.Err{Msg: "failed to lookup url", Status: service.InternalError, Err: dbErr}
		}
	}

	return originalURL, nil
}

func New(conf Config, log logger.Logger, shortenerDB database.URLShortener, keypoolDB database.KeyPool) service.URLShortener {
	return URLShortener{
		urlResultAddress: conf.shortenerAddress,
		shortenerDB:      shortenerDB,
		keyPoolDB:        keypoolDB,
		keyPoolRange:     &database.KeyRange{},
		log:              log,
	}
}
