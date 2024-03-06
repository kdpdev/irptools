package fz

import (
	"errors"
	"fmt"

	"irptools/utils/errs"
)

var (
	ErrPackage = errs.NewPackageError("")

	ErrField                = errs.NewMultiError(errs.NewPackageError("field error"), ErrPackage)
	ErrMissedField          = errs.NewMultiError(errs.NewPackageError("missed field"), ErrField)
	ErrDuplicatedField      = errs.NewMultiError(errs.NewPackageError("duplicated field"), ErrField)
	ErrUnexpectedFieldValue = errs.NewMultiError(errs.NewPackageError("unexpected field value"), ErrField)

	ErrParse   = errs.NewMultiError(errs.NewPackageError("parsing error"), ErrPackage)
	ErrBadLine = errs.NewMultiError(errs.NewPackageError("bad line"), ErrParse)
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewMissedFieldError(field string) error {
	return newFieldError(ErrMissedField, field, "")
}

func NewUnexpectedFieldValueError(field string, value string) error {
	return newFieldError(ErrUnexpectedFieldValue, field, value)
}

func NewDuplicatedFieldError(field string) error {
	return newFieldError(ErrDuplicatedField, field, "")
}

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

func (this *FieldError) Error() string {
	return fmt.Sprintf("%s: field='%s', value = '%s'", this.Head(), this.Field, this.Value)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewBadLineError(line string, details string, cause error) error {
	return errWithCause(newBadLineError(ErrBadLine, line, details), cause)
}

func newBadLineError(head errs.MultiErrorPtr, line string, details string) error {
	return &BadLineError{
		MultiErrorPtr: ensureErrorIs(head, ErrBadLine),
		Line:          line,
		Details:       details,
	}
}

type BadLineError struct {
	errs.MultiErrorPtr
	Line    string
	Details string
}

func (this *BadLineError) Error() string {
	return fmt.Sprintf("%s: line='%s', details:'%s'", this.Head(), this.Line, this.Details)
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
