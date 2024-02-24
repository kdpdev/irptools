package errs

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type contextKeyType int

var (
	contextKeyUnhandledErrorHandler   = contextKeyType(1)
	contextKeyUnhandledErrorDecorator = contextKeyType(2)
)

type UnhandledErrorHandler = func(ctx context.Context, err error)
type UnhandledErrorDecorator = func(ctx context.Context, err error) error

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

func GetUnhandledErrorDecorator(ctx context.Context) UnhandledErrorDecorator {
	return getContextValue(ctx, contextKeyUnhandledErrorDecorator, NoErrorDecorator())
}

func WithUnhandledErrorDecorator(ctx context.Context, decorators ...UnhandledErrorDecorator) context.Context {
	return context.WithValue(ctx, contextKeyUnhandledErrorDecorator, ComposeUnhandledErrorDecorators(decorators...))
}

func getUnhandledErrorHandlerWithoutErrorDecorating(ctx context.Context) UnhandledErrorHandler {
	return getContextValue(ctx, contextKeyUnhandledErrorHandler, NoErrHandler())
}

func GetUnhandledErrorHandler(ctx context.Context) UnhandledErrorHandler {
	h := getContextValue(ctx, contextKeyUnhandledErrorHandler, NoErrHandler())
	return func(ctx context.Context, err error) {
		decorate := GetUnhandledErrorDecorator(ctx)
		h(ctx, decorate(ctx, err))
	}
}

func GetContextedUnhandledErrorHandler(ctx context.Context) func(err error) {
	h := getContextValue(ctx, contextKeyUnhandledErrorHandler, NoErrHandler())
	return func(err error) {
		decorate := GetUnhandledErrorDecorator(ctx)
		h(ctx, decorate(ctx, err))
	}
}

func WithUnhandledErrorHandler(ctx context.Context, h UnhandledErrorHandler) context.Context {
	return context.WithValue(ctx, contextKeyUnhandledErrorHandler, h)
}

func WithUnhandledErrorLogger(ctx context.Context, logf func(format string, args ...interface{})) context.Context {
	return WithUnhandledErrorHandler(ctx, func(ctx context.Context, err error) {
		logf("%v", err)
	})
}

func MakePassIsErrorFilter(errs ...error) func(err error) bool {
	return func(e error) bool {
		for _, err := range errs {
			if errors.Is(e, err) {
				return true
			}
		}
		return false
	}
}

func MakeSkipIsErrorFilter(errs ...error) func(err error) bool {
	pass := MakePassIsErrorFilter(errs...)
	return func(err error) bool {
		return !pass(err)
	}
}

func MakeSkipContextErrorsFilter() func(err error) bool {
	return MakeSkipIsErrorFilter(context.Canceled, context.DeadlineExceeded)
}

func WithUnhandledErrorFilter(ctx context.Context, pass func(err error) bool) context.Context {
	prev := getUnhandledErrorHandlerWithoutErrorDecorating(ctx)
	return WithUnhandledErrorHandler(ctx, func(ctx context.Context, err error) {
		if pass(err) {
			prev(ctx, err)
		}
	})
}

func WithUnhandledErrorContextErrorsFilter(ctx context.Context) context.Context {
	return WithUnhandledErrorFilter(ctx, MakeSkipContextErrorsFilter())
}

func WithUnhandledErrorsCollector(ctx context.Context) (context.Context, func() []error) {
	guard := &sync.RWMutex{}

	errs := make([]error, 0)

	getErrors := func() []error {
		guard.RLock()
		defer guard.RUnlock()
		clone := make([]error, len(errs))
		copy(clone, errs)
		return clone
	}

	prevHandler := getUnhandledErrorHandlerWithoutErrorDecorating(ctx)

	handler := func(ctx context.Context, err error) {
		guard.Lock()
		defer guard.Unlock()
		errs = append(errs, err)
		prevHandler(ctx, err)
	}

	return WithUnhandledErrorHandler(ctx, handler), getErrors
}

func OnUnhandledError(ctx context.Context, err error) {
	GetContextedUnhandledErrorHandler(ctx)(err)
}

func NoErrHandler() UnhandledErrorHandler {
	return func(ctx context.Context, err error) {
	}
}

func NoErrorDecorator() UnhandledErrorDecorator {
	return func(ctx context.Context, err error) error {
		return err
	}
}

func DefaultUnhandledErrorDecorator() UnhandledErrorDecorator {
	return func(ctx context.Context, err error) error {
		return fmt.Errorf("UNHANDLED ERROR: %w", err)
	}
}

func DefaultContextInfoUnhandledErrorDecorator(getCtxInfo func(ctx context.Context) string) UnhandledErrorDecorator {
	return func(ctx context.Context, err error) error {
		return fmt.Errorf("UNHANDLED ERROR: %s%w", getCtxInfo(ctx), err)
	}
}

func ComposeUnhandledErrorDecorators(decorators ...UnhandledErrorDecorator) UnhandledErrorDecorator {
	return func(ctx context.Context, err error) error {
		for _, d := range decorators {
			err = d(ctx, err)
		}
		return err
	}
}
