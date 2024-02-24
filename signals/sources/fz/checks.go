package fz

import (
	"errors"

	"irptools/signals/irp"
	"irptools/utils/alg"
	"irptools/utils/fp"
)

func NewErrorChecker() *ErrorChecker {
	handler := &ErrorChecker{}
	handler.initCheckers()
	return handler
}

type ErrorChecker struct {
	ignoreAllUnsupportedProtocols      bool
	ignoreSpecificUnsupportedProtocols map[string]bool
	checkers                           []fp.FnPred[error]
}

func (this *ErrorChecker) IgnoreAllUnsupportedProtocols(ignoreAll bool) {
	this.ignoreAllUnsupportedProtocols = ignoreAll
}

func (this *ErrorChecker) IgnoreSpecificUnsupportedProtocols(protocols []string) {
	this.ignoreSpecificUnsupportedProtocols = alg.ArrToMap(protocols, true)
}

func (this *ErrorChecker) IsExpectedError(err error) bool {
	if err == nil {
		return false
	}

	for _, check := range this.checkers {
		if check(err) {
			return true
		}
	}

	return false
}

func (this *ErrorChecker) initCheckers() {
	this.checkers = []fp.FnPred[error]{
		this.isExpectedUnsupportedProtocolError,
	}
}

func (this *ErrorChecker) isExpectedUnsupportedProtocolError(err error) bool {
	upe := &irp.UnsupportedProtocolError{}
	if !errors.As(err, &upe) {
		return false
	}

	if this.ignoreAllUnsupportedProtocols {
		return true
	}

	if this.ignoreSpecificUnsupportedProtocols != nil && this.ignoreSpecificUnsupportedProtocols[upe.Protocol] {
		return true
	}

	return false
}
