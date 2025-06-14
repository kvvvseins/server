package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var guidNotFoundError = errors.New("в заголовке не найден Guid пользователя")
var invalidGuidError = errors.New("в заголовке не верный Guid пользователя")

// ErrorResponse структура ошибки сервера.
type ErrorResponse struct {
	Message string `json:"message"`
}

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

// GetUserIDFromRequest получить id пользователя из заголовка
func GetUserIDFromRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	guid := r.Header.Get("X-User-Id")
	if guid == "" {
		ErrorResponseOutput(r.Context(), w, guidNotFoundError, "")

		return uuid.UUID{}, false
	}

	guidParsed, err := uuid.Parse(guid)
	if err != nil {
		ErrorResponseOutput(r.Context(), w, invalidGuidError, "")

		return uuid.UUID{}, false
	}

	return guidParsed, true
}

// ErrorResponseOutput вывод ошибки запроса
func ErrorResponseOutput(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
	errMsg string,
) {
	if errors.Is(err, guidNotFoundError) || errors.Is(err, invalidGuidError) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})

		return
	}

	if err != nil {
		GetLogger(ctx).Error("ошибка запроса", "err", err)
	}

	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: errMsg})
}
