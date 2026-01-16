package service

import "net/url"

type URLKey string

type URLShortener interface {
	CreateShortURL(url.URL, uint) (URLKey, error)
	GetOriginalURL(URLKey) (url.URL, error)
}
