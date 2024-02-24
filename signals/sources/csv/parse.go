package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"irptools/signals/signal"
	"irptools/utils/errs"
)

type Header = string
type FieldSetter = func(s *signal.Signal, value string) error
type Mapping = map[Header]FieldSetter
type SignalConsumer = func(signal signal.Signal) error
type CsvRecord = []string

type FastMappingItem struct {
	Header    Header
	RecordPos int
	Set       FieldSetter
}
type FastMapping = []struct {
	Header    Header
	RecordPos int
	Set       FieldSetter
}

func ParseCsvStream(stream io.Reader, mapping Mapping, consume SignalConsumer) error {
	reader := csv.NewReader(stream)

	headers, err := reader.Read()
	if err != nil {
		return errs.Errorf("failed to read header: %w", err)
	}

	fastMapping, err := makeFastMapping(headers, mapping)
	if err != nil {
		return errs.Wrap(err)
	}

	processRecord := func(r CsvRecord) error {
		var s signal.Signal

		for _, item := range fastMapping {
			err := item.Set(&s, r[item.RecordPos])
			if err != nil {
				return errs.Errorf("failed to set %s field: %w", item.Header, err)
			}
		}

		return consume(s)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errs.Errorf("failed to read csv record: %w", err)
		}

		err = processRecord(record)
		if err != nil {
			return errs.Errorf("failed to process csv record: %w", err)
		}
	}

	return nil
}

func makeFastMapping(headers CsvRecord, mapping Mapping) (FastMapping, error) {
	fastMapping := FastMapping{}
	processed := map[Header]struct{}{}
	for i, h := range headers {
		m, ok := mapping[h]
		if ok {
			fastMapping = append(fastMapping, FastMappingItem{
				Header:    h,
				RecordPos: i,
				Set:       m,
			})

			processed[h] = struct{}{}
		}
	}

	missedHeaders := make([]Header, 0)

	for h := range mapping {
		if _, ok := processed[h]; !ok {
			missedHeaders = append(missedHeaders, h)
		}
	}

	if len(missedHeaders) > 0 {
		return nil, errs.Errorf("missed headers: %v", strings.Join(missedHeaders, ","))
	}

	return fastMapping, nil
}
