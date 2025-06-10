package server

import (
	"context"
	"log/slog"
	"os"
)

type key int

const (
	keyRequestID key = iota
	keyLogger
	keyParentRequestID
)

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(keyRequestID).(string); ok {
		return requestID
	}

	slog.Warn("Not found request_id in context")

	return ""
}

func GetParentRequestID(ctx context.Context) string {
	if parentRequestID, ok := ctx.Value(keyParentRequestID).(string); ok {
		return parentRequestID
	}

	return ""
}

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(keyLogger).(*slog.Logger); ok {
		return logger
	}

	slog.Warn("Not found logger in context")

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
}
