package lambdalog

import "context"

type key struct{}

var contextKey = &key{}

func NewContextWithRecorder(parent context.Context, r *Recorder) context.Context {
	return context.WithValue(parent, contextKey, r)
}

func RecorderFromContext(ctx context.Context) (*Recorder, bool) {
	r, ok := ctx.Value(contextKey).(*Recorder)
	return r, ok
}
