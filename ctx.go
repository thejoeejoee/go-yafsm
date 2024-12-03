package yafasm

import (
	"context"
)

type ctxOriginStateKey struct{}
type ctxTargetStateKey struct{}
type ctxEventKey struct{}
type ctxConditionErr struct{}

func with[T comparable](key any) func(context.Context, T) context.Context {
	return func(ctx context.Context, v T) context.Context {
		//slog.Info("with", "key", key, "v", v)
		return context.WithValue(ctx, key, v)
	}
}

func from[T comparable](key any) func(context.Context) *T {
	return func(ctx context.Context) *T {
		v, ok := ctx.Value(key).(T)
		//slog.Info("from", "key", key, "v", v, "ok", ok)
		if ok {
			return &v
		}
		return nil
	}
}

func withOriginState[S comparable](ctx context.Context, from S) context.Context {
	return with[S](ctxOriginStateKey{})(ctx, from)
}

func withTargetState[S comparable](ctx context.Context, to S) context.Context {
	return with[S](ctxTargetStateKey{})(ctx, to)
}

// withEvent sets the event in the context.
func withEvent[E comparable](ctx context.Context, e E) context.Context {
	return with[E](ctxEventKey{})(ctx, e)
}

// OriginStateFromCtx gets the "from" state from the context.
func OriginStateFromCtx[S comparable](ctx context.Context) *S {
	return from[S](ctxOriginStateKey{})(ctx)
}

// TargetStateFromCtx gets the "to" state from the context.
func TargetStateFromCtx[S comparable](ctx context.Context) *S {
	return from[S](ctxTargetStateKey{})(ctx)
}

// EventFromCtx gets the event from the context.
func EventFromCtx[E comparable](ctx context.Context) *E {
	return from[E](ctxEventKey{})(ctx)
}

// ConditionErrFromCtx gets the error from the context.
func ConditionErrFromCtx(ctx context.Context) *error {
	return from[error](ctxConditionErr{})(ctx)
}

// withConditionErr sets the error in the context.
func withConditionErr(ctx context.Context, err error) context.Context {
	return with[error](ctxConditionErr{})(ctx, err)
}
