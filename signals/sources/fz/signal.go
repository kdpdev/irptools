package fz

import (
	"errors"
	"fmt"
	"strings"

	"irptools/signals/irp"
	"irptools/signals/signal"
	"irptools/utils/errs"
)

func signalFromFields(fields FieldsMap) (signal.Signal, error) {
	s := signal.Signal{}
	err := processField(fields, signalFieldType, func(value string) error {
		decodeSignal, ok := signalDecoders[value]
		if !ok {
			return errs.Wrap(NewUnexpectedFieldValueError(signalFieldType, value))
		}
		var err error
		s, err = decodeSignal(fields)
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

type decodeSignalFn func(fields FieldsMap) (signal.Signal, error)

var signalDecoders = map[string]decodeSignalFn{
	"":       decodeFromUnknownFields,
	"parsed": decodeFromParsedFields,
	"raw":    decodeFromRawFields,
}

func decodeFromUnknownFields(fields FieldsMap) (signal.Signal, error) {
	return signal.Signal{}, errs.Wrap(NewMissedFieldError(signalFieldType))
}

func decodeFromParsedFields(fields FieldsMap) (signal.Signal, error) {
	s, err := decodeCommonFields(fields)
	if err != nil {
		return s, errs.Wrap(err)
	}

	var code irp.SignalCode
	err = processField(fields, signalFieldAddress, func(value string) error {
		code.Address, err = irp.ParseHex32(value)
		return errs.Wrap(err)
	})
	if err != nil {
		return s, err
	}

	err = processField(fields, signalFieldCommand, func(value string) error {
		code.Command, err = irp.ParseHex32(value)
		return errs.Wrap(err)
	})
	if err != nil {
		return s, err
	}

	err = processField(fields, signalFieldProtocol, func(protocol string) error {
		irpDecoder, err := irp.GetIrp(protocol)
		if err == nil {
			s.Data, err = irpDecoder.Decode(code)
		}
		return errs.Wrap(err)
	})

	return s, err
}

func decodeFromRawFields(fields FieldsMap) (signal.Signal, error) {
	s, err := decodeCommonFields(fields)
	if err != nil {
		return s, errs.Wrap(err)
	}

	s.Protocol = "raw"

	err = processField(fields, signalFieldData, func(data string) error {
		s.Data, err = irp.SplitToMicrosArr(data, " ")
		return errs.Wrap(err)
	})

	return s, err
}

func decodeCommonFields(fields FieldsMap) (signal.Signal, error) {
	freq, err := getFrequency(fields)
	if err != nil {
		return signal.Signal{}, errs.Wrap(err)
	}

	s := signal.Signal{}
	s.Function = fields[signalFieldName]
	s.Protocol = fields[signalFieldProtocol]
	s.Device = "unknown"
	s.Frequency = freq

	return s, nil
}

func getFrequency(fields FieldsMap) (irp.Frequency, error) {
	freqStr, ok := fields[signalFieldFrequency]
	if ok {
		freq, err := irp.ParseFrequency(freqStr)
		if err != nil {
			return 0, errs.Wrap(NewFieldError(signalFieldFrequency, freqStr, err))
		}
		return freq, nil
	}

	freq := irp.Frequency(0)
	err := processField(fields, signalFieldProtocol, func(protocol string) error {
		irpFreq, err := irp.GetIrp(protocol)
		if err == nil {
			freq = irpFreq.Frequency()
		}
		return errs.Wrap(err)
	})

	return freq, err
}

func processField(fields FieldsMap, field string, exec func(value string) error) error {
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
