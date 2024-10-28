package logger

import (
	"context"
	"log/slog"
	"testing"
)

func TestLogger(t *testing.T) {
	Init()
	ctx := AppendCtx(context.Background(), slog.String("request_id", "req-123"))
	id := GenerateTraceID()
	ctx = AppendCtx(ctx, slog.String(TraceID, id))

	slog.InfoContext(ctx, "image uploaded", slog.String("image_id", "img-998"))
	traceID := GetTraceID(ctx)
	if traceID != id {
		t.Errorf("expected trace id %s, got %s", id, traceID)
	}
}
