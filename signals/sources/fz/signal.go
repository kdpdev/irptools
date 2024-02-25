package fz

import (
	"errors"
	"fmt"
	"strings"

	"irptools/signals/irp"
	"irptools/signals/signal"
	"irptools/utils/errs"
)

func signalFromFields(fields FieldsMap, checker *ErrorChecker) (signal.Signal, error) {
	s := signal.Signal{}
	d := decoder{checker: checker}

	err := d.processField(fields, signalFieldType, func(value string) error {
		decodeSignal, ok := signalDecoders[value]
		if !ok {
			return errs.Wrap(NewUnexpectedFieldValueError(signalFieldType, value))
		}
		var err error
		s, err = decodeSignal(&d, fields)
		return errs.Wrap(err)
	})

	if err != nil {
		b := strings.Builder{}
		b.WriteString("failed to build signal from batch:\n")
		for k, v := range fields {
			_, _ = fmt.Fprintf(&b, "  %s: '%s'\n", k, v)
		}
		msg := b.String()
		msg = msg[:len(msg)-1]
		return s, errs.Join(ErrPackage, errors.New(msg), err)
	}

	return s, errs.Wrap(err)
}

type decodeSignalFn func(d *decoder, fields FieldsMap) (signal.Signal, error)

var signalDecoders = map[string]decodeSignalFn{
	"":       func(d *decoder, fields FieldsMap) (signal.Signal, error) { return d.decodeFromUnknownFields(fields) },
	"parsed": func(d *decoder, fields FieldsMap) (signal.Signal, error) { return d.decodeFromParsedFields(fields) },
	"raw":    func(d *decoder, fields FieldsMap) (signal.Signal, error) { return d.decodeFromRawFields(fields) },
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type checkingIrp struct {
	checker *ErrorChecker
	origin  irp.Irp
}

func (this checkingIrp) Protocol() string {
	return this.origin.Protocol()
}

func (this checkingIrp) Frequency() irp.Frequency {
	return this.origin.Frequency()
}

func (this checkingIrp) Decode(code irp.SignalCode) (irp.SignalData, error) {
	data, err := this.origin.Decode(code)
	if err == nil {
		return data, nil
	}

	if this.checker.IsExpectedError(err) {
		return nil, nil
	}

	return nil, errs.Wrap(err)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type decoder struct {
	checker *ErrorChecker
}

func (this *decoder) getIrp(protocol string) (irp.Irp, error) {
	originIrp, err := irp.GetIrp(protocol)
	if err == nil {
		return checkingIrp{checker: this.checker, origin: originIrp}, nil
	}
	if this.checker.IsExpectedError(err) {
		return checkingIrp{checker: this.checker, origin: irp.NewIrpUnsupported(protocol, 0)}, nil
	}
	return nil, errs.Wrap(err)
}

func (this *decoder) decodeFromUnknownFields(fields FieldsMap) (signal.Signal, error) {
	return signal.Signal{}, errs.Wrap(NewMissedFieldError(signalFieldType))
}

func (this *decoder) decodeFromParsedFields(fields FieldsMap) (signal.Signal, error) {
	s := signal.Signal{}
	s.Function = fields[signalFieldName]

	err := this.processField(fields, signalFieldAddress, func(value string) error {
		var err error
		s.Code.Address, err = irp.ParseHex32(value)
		return errs.Wrap(err)
	})
	if err != nil {
		return s, err
	}

	err = this.processField(fields, signalFieldCommand, func(value string) error {
		s.Code.Command, err = irp.ParseHex32(value)
		return errs.Wrap(err)
	})
	if err != nil {
		return s, err
	}

	err = this.processField(fields, signalFieldProtocol, func(protocol string) error {
		s.Protocol = protocol
		irp, err := this.getIrp(protocol)
		if err != nil {
			return errs.Wrap(err)
		}
		s.Frequency = irp.Frequency()
		s.Data, err = irp.Decode(s.Code)
		return errs.Wrap(err)
	})

	return s, err
}

func (this *decoder) decodeFromRawFields(fields FieldsMap) (signal.Signal, error) {
	s := signal.Signal{}
	s.Function = fields[signalFieldName]

	err := this.processField(fields, signalFieldFrequency, func(data string) error {
		freq, err := irp.ParseFrequency(data)
		if err != nil {
			return errs.Wrap(err)
		}
		s.Frequency = freq
		return nil
	})

	err = this.processField(fields, signalFieldData, func(data string) error {
		s.Data, err = irp.SplitToMicrosArr(data, " ")
		return errs.Wrap(err)
	})

	return s, err
}

func (this *decoder) processField(fields FieldsMap, field string, exec func(value string) error) error {
	value, ok := fields[field]
	if !ok {
		return errs.Wrap(NewMissedFieldError(field))
	}

	err := exec(value)
	if err != nil {
		return errs.Wrap(NewFieldError(field, value, err))
	}

	return nil
}

const (
	signalFieldProtocol  = "protocol"
	signalFieldFrequency = "frequency"
	signalFieldName      = "name"
	signalFieldType      = "type"
	signalFieldCommand   = "command"
	signalFieldAddress   = "address"
	signalFieldData      = "data"
)
