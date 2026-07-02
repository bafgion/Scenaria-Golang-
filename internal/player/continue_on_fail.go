package player

import "context"

type continueOnFailKey struct{}

// WithContinueOnFail sets whether the runner keeps executing after a scenario failure.
func WithContinueOnFail(ctx context.Context, enabled bool) context.Context {
	if !enabled {
		return ctx
	}
	return context.WithValue(ctx, continueOnFailKey{}, true)
}

// ContinueOnFail reports whether failed scenarios should not abort the rest of the suite.
func ContinueOnFail(ctx context.Context) bool {
	v, _ := ctx.Value(continueOnFailKey{}).(bool)
	return v
}
