package service

import "net/url"

type URLKey string

type URLShortener interface {
	CreateShortURL(*url.URL) (URLKey, error)
	GetOriginalURL(URLKey) (*url.URL, error)
}
