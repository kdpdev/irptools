package fp

type FnPred[T any] func(T) bool
type FnVoid[T any] func(T)
type FnR[T, R any] func(T) R
type FnE[T any] FnR[T, error]
type FnRE[T, R any] func(T) (R, error)
type FnResErr[T, R any] func(T) ResErr[R]

type ResErr[T any] struct {
	Res T
	Err error
}

func Not[T any](pred FnPred[T]) FnPred[T] {
	return func(t T) bool {
		return !pred(t)
	}
}

func Flow[T, R any](from FnRE[T, R], next FnE[R]) FnE[T] {
	return func(t T) error {
		res, err := from(t)
		if err != nil {
			return err
		}
		return next(res)
	}
}

func FlowFilterError[T, R any](
	from FnRE[T, R],
	pred FnPred[error],
	next FnE[R]) FnE[T] {

	return func(t T) error {
		res, err := from(t)
		if err != nil {
			if pred(err) {
				return nil
			}
			return err
		}
		return next(res)
	}
}

func Filter[T any](pred FnPred[T], next FnE[T]) FnE[T] {
	return func(t T) error {
		if pred(t) {
			return next(t)
		}
		return nil
	}
}

func Split[T any](
	isSplitter FnPred[T],
	includeSplitter bool,
	ignoreUntilSplitter bool,
	next FnE[[]T]) FnE[T] {

	batch := []T{}

	process := func(data T) error {
		if isSplitter(data) {
			err := next(batch)
			if err != nil {
				return err
			}
			batch = []T{}
			if includeSplitter {
				batch = append(batch, data)
			}
			return nil
		}
		batch = append(batch, data)
		return nil
	}

	if ignoreUntilSplitter {
		originProcess := process
		process = func(data T) error {
			if isSplitter(data) {
				if includeSplitter {
					batch = append(batch, data)
				}
				process = originProcess
				return nil
			}
			return nil
		}
	}

	return func(t T) error {
		return process(t)
	}
}

func FnVoid2E[T any](f FnVoid[T]) FnE[T] {
	return func(t T) error {
		f(t)
		return nil
	}
}

func FnE2Void[T any](f FnE[T]) FnVoid[T] {
	return func(t T) {
		_ = f(t)
	}

}

func FnR2RE[T, R any](f FnR[T, R]) FnRE[T, R] {
	return func(t T) (R, error) {
		r := f(t)
		return r, nil
	}
}

func FnRE2ResErr[T, R any](f FnRE[T, R]) FnResErr[T, R] {
	return func(t T) ResErr[R] {
		r, e := f(t)
		return ResErr[R]{
			Res: r,
			Err: e,
		}
	}
}
