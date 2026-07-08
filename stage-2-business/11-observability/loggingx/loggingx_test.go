package loggingx

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"just-go/stage-2-business/11-observability/tracex"
)

func TestInfoWritesTraceAwareJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, slog.LevelInfo)
	ctx := tracex.WithTraceID(context.Background(), "trace-log")

	Info(ctx, logger, "article.created", slog.Int("article_id", 42))

	var record map[string]any
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("log should be json: %v\n%s", err, buf.String())
	}
	if record["msg"] != "article.created" || record["trace_id"] != "trace-log" || record["article_id"] != float64(42) {
		t.Fatalf("unexpected log record: %#v", record)
	}
}
