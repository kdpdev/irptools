package errs

import "fmt"

func makeRequiredErrorWithSkip(skip int, name string) error {
	return newCheckErrWithSkipf(2+skip, "%v is required", name)
}

func makeNotValidErrorWithSkip(skip int, err error, name string) error {
	return newCheckErrWithSkipf(2+skip, "%v is not valid: %v", name, err)
}

func makeNegativeErrorWithSkip(skip int, name string) error {
	return newCheckErrWithSkipf(2+skip, "%v is negative", name)
}

func makeNotPositiveErrorWithSkip(skip int, name string) error {
	return newCheckErrWithSkipf(2+skip, "%v is <= 0", name)
}

func makeNotZeroErrorWithSkip(skip int, name string) error {
	return newCheckErrWithSkipf(2+skip, "%v is not zero", name)
}

func makeZeroErrorWithSkip(skip int, name string) error {
	return newCheckErrWithSkipf(2+skip, "%v is zero", name)
}

func newCheckErrWithSkipf(skip int, format string, args ...any) error {
	ensureWithoutW(format)
	return WithFrameSkip(fmt.Errorf(format, args...), skip)
}
