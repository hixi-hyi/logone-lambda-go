package lambdalog

import "context"

type key struct{}

var contextKey = &key{}

func NewContextWithLogger(parent context.Context, r *Logger) context.Context {
	return context.WithValue(parent, contextKey, r)
}

func LoggerFromContext(ctx context.Context) (*Logger, bool) {
	r, ok := ctx.Value(contextKey).(*Logger)
	return r, ok
}
