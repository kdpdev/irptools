package errs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var errJoined = errors.New("joined error")

func Join(errs ...error) error {
	joined := errors.Join(errs...)
	if joined == nil {
		return nil
	}

	wrapped, ok := joined.(interface{ Unwrap() []error })
	if !ok {
		panic("unexpected implementation of errors.Join(errs ...error) error")
	}

	errs = wrapped.Unwrap()
	if len(errs) == 0 {
		return joined
	}

	if len(errs) == 1 {
		return errs[0]
	}

	return &joinError{errs: append([]error{errJoined}, errs...)}
}

type joinError struct {
	errs []error
}

func (this *joinError) Format(s fmt.State, verb rune) {
	this.print(s, verb, nil)
}

func (this *joinError) Error() string {
	var b strings.Builder
	this.print(&b, 's', func(e error) interface{} { return e.Error() })
	return b.String()
}

func (this *joinError) print(w io.Writer, verb rune, errData func(e error) interface{}) {
	verbFmt := string([]rune{'\n', 'e', '[', '%', 'd', '/', '%', 'd', ']', ':', ' ', '%', verb})
	nested := nestedWriter{w}
	_, _ = fmt.Fprintf(&nested, verbFmt[len(verbFmt)-2:], this.errs[0])
	count := len(this.errs) - 1
	for i, err := range this.errs[1:] {
		var data interface{} = err
		if errData != nil {
			data = errData(err)
		}
		_, _ = fmt.Fprintf(&nested, verbFmt, i+1, count, data)
	}
}

func (e *joinError) Unwrap() []error {
	return e.errs
}

type nestedWriter struct {
	io.Writer
}

func (this *nestedWriter) Write(b []byte) (n int, err error) {
	write := func(data []byte) bool {
		nn, e := this.Writer.Write(data)
		if e != nil {
			err = e
			return false
		}
		n += nn
		return true
	}

	lines := bytes.Split(b, []byte{'\n'})
	indent := [...]byte{'\n', ' ', ' '}

	for _, l := range lines[:len(lines)-1] {
		if !write(l) {
			return
		}
		if !write(indent[:]) {
			return
		}
	}

	write(lines[len(lines)-1])
	return
}
