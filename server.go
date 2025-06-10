package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// NewServer создает новый инстанс роутера и сервера.
func NewServer(port int, logLevel slog.Level) (*chi.Mux, *http.Server) {
	router := chi.NewRouter()

	router.Use(
		Recovery,
		LoggingHandler(logLevel),
	)

	accessLog := AccessLog(zerolog.New(os.Stdout))

	server := &http.Server{
		Handler:           accessLog(router),
		Addr:              fmt.Sprintf(":%d", port),
		ReadTimeout:       time.Millisecond * 250,
		ReadHeaderTimeout: time.Millisecond * 200,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Minute * 30,
	}

	return router, server
}

// AddRequestIDToRequestHeader добавляет request id в заголовки и родительский request id.
func AddRequestIDToRequestHeader(header http.Header, parentRequestID string) {
	header.Set(traceParentHeader, parentRequestID)
	header.Set(requestIdHeader, uuid.NewString())
}
