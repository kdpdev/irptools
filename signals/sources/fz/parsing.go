package fz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"irptools/signals/signal"
	"irptools/utils/errs"
	"irptools/utils/fp"
	"irptools/utils/misc"
)

type parseCfg struct {
	brand                              string
	source                             string
	ignoreAllUnsupportedProtocols      bool
	ignoreSpecificUnsupportedProtocols []string
}

func ParseIrStream(
	cfg parseCfg,
	stream io.Reader,
	consumer SignalConsumer) (int, error) {

	errorsChecker := NewErrorChecker()
	errorsChecker.IgnoreAllUnsupportedProtocols(cfg.ignoreAllUnsupportedProtocols)
	errorsChecker.IgnoreSpecificUnsupportedProtocols(cfg.ignoreSpecificUnsupportedProtocols)

	nextId := func() func() string {
		id := -1
		return func() string {
			id++
			return strconv.Itoa(id)
		}
	}()

	enrichSignal := func(s signal.Signal) signal.Signal {
		s.Id = nextId()
		s.Source = cfg.source
		s.Brand = cfg.brand
		return s
	}

	signalsCount := 0

	consume := func(signal signal.Signal) error {
		err := consumer.Consume(signal)
		if err != nil {
			return err
		}
		signalsCount++
		return nil
	}

	fieldsToSignal := func(fields FieldsMap) (signal.Signal, error) {
		return signalFromFields(fields, errorsChecker)
	}

	signals := fp.FlowFilterError(fieldsToSignal,
		errorsChecker.IsExpectedError,
		fp.Flow(fp.FnR2RE(enrichSignal), consume))

	linesToBatch := fp.Split(isLinesBatchSplitter, true, true,
		fp.Filter(fp.Not(isLinesBatchEmpty),
			fp.Flow(fieldsFromLinesBatch, signals)))

	lines := fp.Filter(isStreamLineValid,
		fp.Flow(fp.FnR2RE(strings.TrimSpace),
			linesToBatch))

	stream = io.MultiReader(stream, strings.NewReader("\n"+signalsEndOfLines))
	err := misc.EnumStreamLines(io.NopCloser(stream), func(line string) (bool, error) {
		return true, lines(line)
	})

	return signalsCount, errs.Wrap(err)
}

func isLinesBatchEmpty(batch LinesBatch) bool {
	return len(batch) == 0
}

func fieldsFromLinesBatch(lines LinesBatch) (FieldsMap, error) {
	parseLine := func(line string) (string, string, error) {
		const separator = ":"
		separatorPos := strings.Index(line, separator)
		if separatorPos < 0 {
			return "", "", errs.Wrap(NewBadLineError(line, fmt.Sprintf("line is without the '%s' separator", separator), nil))
		}
		key := strings.ToLower(strings.TrimSpace(line[0:separatorPos]))
		value := strings.TrimSpace(line[separatorPos+1:])
		return key, value, nil
	}

	fields := FieldsMap{}

	processLine := func(line string) error {
		key, value, err := parseLine(line)
		if err != nil {
			return errs.Wrap(err)
		}

		if _, ok := fields[key]; ok {
			return errs.Wrap(NewDuplicatedFieldError(key))
		}

		fields[key] = value

		return nil
	}

	makeError := func(err error) error {
		msg := strings.Builder{}
		msg.WriteString("failed to parse signal lines:\n")
		for i, line := range lines {
			msg.WriteString(fmt.Sprintf("  %d: '%s'\n", i, line))
		}
		return fmt.Errorf("%s: %w", msg.String(), err)
	}

	for _, line := range lines {
		err := processLine(line)
		if err != nil {
			return nil, errs.Wrap(makeError(err))
		}
	}

	return fields, nil
}

func isStreamLineValid(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}
	if line[0] == '#' {
		return false
	}
	return true
}

func isLinesBatchSplitter(line string) bool {
	return strings.Index(line, signalFieldName) == 0 ||
		line == signalsEndOfLines
}

const signalsEndOfLines = "!!!EOL!!!"
