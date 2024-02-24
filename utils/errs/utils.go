package errs

import (
	"reflect"
)

type Validatable interface {
	Validate() error
}

func Catch(exec func()) (result error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		err, ok := r.(error)
		if ok {
			result = err
			return
		}

		result = Errorf("panic: %v", r)
	}()

	exec()

	return
}

func Throw(obj any) {
	panic(obj)
}

func ThrowIf(obj any) {
	if obj != nil {
		Throw(obj)
	}
}

func isNil(obj any) bool {
	if obj == nil {
		return true
	}

	value := reflect.ValueOf(obj)
	kind := value.Kind()
	switch kind {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Ptr,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		return value.IsNil()
	}

	return false
}
