package logx

import (
	"context"
)

type contextKey string

func contextWithIDs(ctx context.Context, traceID, spanID string) context.Context {
	ctx = context.WithValue(ctx, contextKey("trace_id"), traceID)
	ctx = context.WithValue(ctx, contextKey("span_id"), spanID)
	return ctx
}

func ContextWithTraceSpan(ctx context.Context, traceID, spanID string) context.Context {
	return contextWithIDs(ctx, traceID, spanID)
}