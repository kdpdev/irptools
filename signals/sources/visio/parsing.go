package visio

import (
	"io"
	"strings"

	"irptools/signals/irp"
	"irptools/signals/signal"
	"irptools/signals/sources/csv"
	"irptools/utils/errs"
)

type parseCfg struct {
	source string
}

type CsvRecord = []string

func ParseCsvStream(
	cfg parseCfg,
	stream io.Reader,
	consumer SignalConsumer) (int, error) {

	count := 0
	err := csv.ParseCsvStream(stream, csvMapping, func(signal signal.Signal) error {
		signal.Source = cfg.source
		err := consumer.Consume(signal)
		if err != nil {
			return errs.Wrap(err)
		}
		count++
		return nil
	})

	return count, errs.Wrap(err)
}

var csvMapping = csv.Mapping{
	"_id":             withFieldError("Id", setSignalId),
	"brand":           withFieldError("Brand", setSignalBrand),
	"device":          withFieldError("Device", setSignalDevice),
	"button_fragment": withFieldError("Function", setSignalFunction),
	"frequency":       withFieldError("Frequency", setSignalFrequency),
	"main_frame":      withFieldError("Data", setSignalData),
}

func withFieldError(field string, setter csv.FieldSetter) csv.FieldSetter {
	return func(s *signal.Signal, value string) error {
		err := setter(s, value)
		if err != nil {
			return NewFieldError(field, value, err)
		}
		return nil
	}
}

func setSignalId(s *signal.Signal, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errs.Error("empty id")
	}
	return setStringField(&s.Id, value)
}

func setSignalBrand(s *signal.Signal, value string) error {
	return setStringField(&s.Brand, value)
}

func setSignalDevice(s *signal.Signal, value string) error {
	return setStringField(&s.Device, value)
}

func setSignalFunction(s *signal.Signal, value string) error {
	return setStringField(&s.Function, value)
}

func setStringField(field *string, value string) error {
	*field = strings.TrimSpace(value)
	return nil
}

func setSignalFrequency(s *signal.Signal, value string) error {
	freq, err := irp.ParseFrequency(value)
	if err != nil {
		return errs.Wrap(err)
	}
	s.Frequency = freq
	return nil
}

func setSignalData(s *signal.Signal, value string) error {
	data, err := irp.SplitToMicrosArr(value, ",")
	if err != nil {
		return errs.Wrap(err)
	}
	s.Data = data
	return nil
}
