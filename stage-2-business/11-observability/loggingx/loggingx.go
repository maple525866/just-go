package loggingx

import (
	"context"
	"io"
	"log/slog"

	"just-go/stage-2-business/11-observability/tracex"
)

func NewJSONLogger(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))
}

func Info(ctx context.Context, logger *slog.Logger, msg string, attrs ...slog.Attr) {
	fields := make([]any, 0, len(attrs)+1)
	if traceID := tracex.TraceID(ctx); traceID != "" {
		fields = append(fields, slog.String("trace_id", traceID))
	}
	for _, attr := range attrs {
		fields = append(fields, attr)
	}
	logger.InfoContext(ctx, msg, fields...)
}
