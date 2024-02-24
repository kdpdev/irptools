package visio

import (
	"errors"
	"fmt"

	"irptools/utils/errs"
)

var (
	ErrPackage = errs.NewPackageError("")

	ErrField = errs.NewMultiError(errs.NewPackageError("field error"), ErrPackage)
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewFieldError(field string, value string, cause error) error {
	return errWithCause(newFieldError(nil, field, value), cause)
}

func newFieldError(head errs.MultiErrorPtr, field, value string) error {
	return &FieldError{
		MultiErrorPtr: ensureErrorIs(head, ErrField),
		Field:         field,
		Value:         value,
	}
}

type FieldError struct {
	errs.MultiErrorPtr
	Field string
	Value string
}

func (this FieldError) Error() string {
	return fmt.Sprintf("%s: field='%s', value = '%s'", this.Head(), this.Field, this.Value)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ensureErrorIs(err errs.MultiErrorPtr, is errs.MultiErrorPtr) errs.MultiErrorPtr {
	if err == nil {
		err = is
	}

	if !errors.Is(err, is) {
		panic(fmt.Sprintf("err is not %v", is.Error()))
	}

	return err
}

func errWithCause(err error, cause error) error {
	return errs.Join(err, cause)
}
