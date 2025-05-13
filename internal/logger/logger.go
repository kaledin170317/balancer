package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

const contextKey = "logger-context-key"

var (
	once sync.Once
	base *slog.Logger
)

func Init(level slog.Leveler) {
	once.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
		base = slog.New(handler)
	})
}

func Base() *slog.Logger {
	if base == nil {
		Init(slog.LevelInfo)
	}
	return base
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Base()
	}
	if l, ok := ctx.Value(contextKey).(*slog.Logger); ok {
		return l
	}
	return Base()
}

func Info(ctx context.Context, msg string, args ...any)  { FromContext(ctx).Info(msg, args...) }
func Debug(ctx context.Context, msg string, args ...any) { FromContext(ctx).Debug(msg, args...) }
func Warn(ctx context.Context, msg string, args ...any)  { FromContext(ctx).Warn(msg, args...) }
func Error(ctx context.Context, msg string, args ...any) { FromContext(ctx).Error(msg, args...) }
