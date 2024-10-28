package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type ContextHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		traceID := span.SpanContext().TraceID().String()
		r.AddAttrs(slog.String(TraceID, traceID))
		return h.Handler.Handle(ctx, r)
	}

	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

// Get TraceID from context
func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, a := range v {
			if a.Key == TraceID {
				return a.Value.String()
			}
		}
	}

	return ""
}

func Init() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			// 獲取檔案名稱跟上層資料夾名稱
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = filepath.Join(filepath.Base(filepath.Dir(s.File)), filepath.Base(s.File))
			}

			if a.Key == slog.TimeKey {
				a = slog.String("time", time.Now().Format("2006/01/02 15:04:05"))
			}
			return a
		},
	})

	h := &ContextHandler{handler}

	logger := slog.New(h)
	slog.SetDefault(logger)
}
