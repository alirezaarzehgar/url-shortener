package slogger

import (
	"log/slog"
	"os"

	"github.com/alirezaarzehgar/url-shortener/internal/logger"
)

type slogger struct {
	logger *slog.Logger
}

func (l slogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l slogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l slogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l slogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func New() logger.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	return slogger{
		logger: slog.New(handler),
	}
}
