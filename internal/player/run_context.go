package player

import "context"

type runIDCtxKey struct{}

// WithRunID attaches a run identifier to the execution context.
func WithRunID(ctx context.Context, runID string) context.Context {
	if runID == "" {
		return ctx
	}
	return context.WithValue(ctx, runIDCtxKey{}, runID)
}

func RunIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(runIDCtxKey{}).(string)
	return id
}
