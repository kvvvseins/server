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

// GetLogger Получить логер
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(keyLogger).(*slog.Logger); ok {
		return logger
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Warn("Not found logger in context")

	return logger
}
