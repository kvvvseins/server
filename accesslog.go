package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/mash/go-accesslog"
	"github.com/rs/zerolog"
)

const (
	traceParentHeader = "traceparent"
	rfc3339Milli      = "2006-01-02T15:04:05.999Z07:00"
	requestIdHeader   = "X-Request-Id"
)

// AccessLog middleware для формирования access-логов в нужном формате.
func AccessLog(l zerolog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logger := &accessLogger{log: l}

		handler := accesslog.NewLoggingHandler(next, logger)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIdHeader)
			if err := uuid.Validate(requestID); err != nil {
				http.Error(w, "Заголовок "+requestIdHeader+" обязателен (UUID)", http.StatusBadRequest)

				return
			}

			parentRequestID := r.Header.Get(traceParentHeader)
			if err := uuid.Validate(parentRequestID); err != nil {
				parentRequestID = ""
				r.Header.Add(traceParentHeader, parentRequestID)
			}

			if requestID == parentRequestID {
				http.Error(w, "Заголовки "+requestIdHeader+" и "+traceParentHeader+" должны отличаться", http.StatusBadRequest)

				return
			}

			ctx := context.WithValue(r.Context(), keyRequestID, requestID)
			ctx = context.WithValue(ctx, keyParentRequestID, parentRequestID)

			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type accessLogger struct {
	log zerolog.Logger
}

func (a accessLogger) Log(record accesslog.LogRecord) {
	if record.Uri == "/favicon.ico" {
		return
	}

	a.log.Log().
		Str("date_time", record.Time.Format(rfc3339Milli)).
		Int64("bytes_sent", record.Size).
		Dur("request_time", record.ElapsedTime).
		Str("host", record.Host).
		Str("request_query", record.Uri).
		Str("http_method", record.Method).
		Int("http_status", record.Status).
		Str("traceparent", record.RequestHeader.Get(traceParentHeader)).
		Str("request_id", record.RequestHeader.Get(requestIdHeader)).
		Send()
}
