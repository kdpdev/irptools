package logs

import (
	"context"
	"fmt"
)

func L(ctx context.Context) Logger {
	return getContextValue(ctx, contextKeyLogger, DefaultLogger())
}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

func WithPrefixedLogger(ctx context.Context, prefix string) (context.Context, Logger) {
	prefixed := L(ctx).WithPrefix(prefix)
	return WithLogger(ctx, prefixed), prefixed
}

func GetScope(ctx context.Context) string {
	return getContextValue(ctx, contextKeyScope, "")
}

func WithScope(ctx context.Context, scope string) (context.Context, Logger) {
	newScope := fmt.Sprintf("%s%s: ", GetScope(ctx), scope)
	ctx = context.WithValue(ctx, contextKeyScope, newScope)
	return WithPrefixedLogger(ctx, scope)
}

func WithCallerScope(ctx context.Context) (context.Context, Logger) {
	return WithScope(ctx, callerFuncName())
}

func getContextValue[T any](ctx context.Context, key contextKeyType, defaultValue T) T {
	val := ctx.Value(key)
	if val == nil {
		return defaultValue
	}
	typedVal, ok := val.(T)
	if !ok {
		return defaultValue
	}
	return typedVal
}

type contextKeyType int

var (
	contextKeyScope  = contextKeyType(1)
	contextKeyLogger = contextKeyType(2)
)
