package http

import (
	"net/http"

	"github.com/alirezaarzehgar/url-shortener/internal/service"
	"github.com/alirezaarzehgar/url-shortener/internal/transport"
)

type HttpTransport struct {
	svc     service.URLShortener
	address string
}

func (t HttpTransport) Run() error {
	mux := http.NewServeMux()

	t.router(mux)

	srv := http.Server{
		Addr:    t.address,
		Handler: mux,
	}

	return srv.ListenAndServe()
}

func (t HttpTransport) router(mux *http.ServeMux) {
	mux.HandleFunc("/", t.controllerRedirectToOriginalURL)
	mux.HandleFunc("/new", t.controllerCreateNewShortLinkFromOriginalURL)
}

func New(conf Config, service service.URLShortener) transport.URLShortener {
	return HttpTransport{
		address: conf.shortenerAddress,
		svc:     service,
	}
}
