package irp

import (
	"fmt"

	"irptools/utils/errs"
)

var (
	Wrap   = errs.Wrap
	Error  = errs.Error
	Errorf = errs.Errorf
)

var (
	ErrPackage             = errs.NewPackageError("")
	ErrUnsupportedProtocol = errs.NewMultiError(errs.NewPackageError("unsupported protocol"), ErrPackage)
)

func NewUnsupportedProtocolError(protocol string) *UnsupportedProtocolError {
	return &UnsupportedProtocolError{
		MultiErrorPtr: ErrUnsupportedProtocol,
		Protocol:      protocol,
	}
}

type UnsupportedProtocolError struct {
	errs.MultiErrorPtr
	Protocol string
}

func (this *UnsupportedProtocolError) Error() string {
	return fmt.Sprintf("%s: %s", this.Head(), this.Protocol)
}
