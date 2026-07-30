package api

import (
	"net/http"
	"time"
)

// NewServer создаёт и настраивает HTTP сервер.
func NewServer(addr string, handler *Handler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", handler.Convert)
	mux.HandleFunc("/currencies", handler.Currencies)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
}
