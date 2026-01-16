package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/alirezaarzehgar/url-shortener/internal/service"
)

const (
	MinURLKeyLen = 4
	MaxURLKeyLen = 12

	TTLMax     uint = 60 * 60 * 24 * 30 * 6
	TTLDefault uint = 604800
)

func (t HttpTransport) controllerCreateNewShortLinkFromOriginalURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
		TTL uint   `json:"ttl"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL, err := url.Parse(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.TTL > TTLMax {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.TTL == 0 {
		req.TTL = TTLDefault
	}

	shortURL, err := t.svc.CreateShortURL(*originalURL, req.TTL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"url":"` + shortURL + `"}`))
}

func (t HttpTransport) controllerRedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if strings.Count(r.URL.Path, "/") != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortKey := r.URL.Path[1:]
	if len(shortKey) < MinURLKeyLen || len(shortKey) > MaxURLKeyLen {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := t.svc.GetOriginalURL(service.URLKey(shortKey))
	if err != nil {
		if svcErr := err.(service.Err); svcErr.NotFound() {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, shortURL.String(), http.StatusPermanentRedirect)
}
