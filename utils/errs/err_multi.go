package errs

import (
	"errors"
)

func NewMultiError(head error, others ...error) MultiErrorPtr {
	arr := make([]error, 1+len(others))
	arr[0] = head
	copy(arr[1:], others)
	return &multiError{
		head: head,
		all:  Join(arr...),
	}
}

type MultiErrorPtr = *multiError

type multiError struct {
	head error
	all  error
}

func (this *multiError) Is(err error) bool {
	return this == err || errors.Is(this.all, err)
}

func (this *multiError) Error() string {
	return this.Head().Error()
	//return this.joined.Error()
}

func (this *multiError) Head() error {
	target := &multiError{}
	ok := errors.As(this.head, &target)
	if ok {
		return target.Head()
	}
	return this.head
}

func (this *multiError) All() error {
	return this.all
}
