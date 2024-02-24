package visio

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"irptools/signals/signal"
	"irptools/utils/errs"
	"irptools/utils/fs"
	jsonutils "irptools/utils/json"
	"irptools/utils/logs"
)

type Options struct {
}

type SignalConsumer interface {
	Consume(signal signal.Signal) error
}

type ClosableSignalConsumer interface {
	SignalConsumer
	io.Closer
}

type SignalConsumerSource = func(filePath string) (ClosableSignalConsumer, error)

func ParseCsvFiles(
	ctx context.Context,
	rootPath string,
	opts interface{},
	getConsumer SignalConsumerSource) (parsedSignalsCount int, err error) {

	l := logs.L(ctx)

	options := Options{}
	err = jsonutils.Cast(opts, &options)
	if err != nil {
		return 0, errs.Wrap(err)
	}

	err = fs.EnumFilePathsWithExt(rootPath, ".csv", func(filePath string) (res bool, err error) {
		dirPath, fileName := filepath.Split(filePath)
		sConsumer := &splittingConsumer{
			fileName:    fileName,
			dirPath:     dirPath,
			getConsumer: getConsumer,
			consumers:   map[string]map[string]ClosableSignalConsumer{},
		}

		defer func() {
			closeErr := sConsumer.Close()
			if err == nil {
				err = closeErr
			}
		}()

		count, err := ParseCsvFile(filePath, options, sConsumer)

		parsedSignalsCount += count
		if err != nil {
			l.I("ERR: %-4d: %s", count, filePath)
		} else {
			l.I(" OK: %-4d: %s", count, filePath)
		}

		if err != nil {
			return false, errs.Errorf("failed to parse '%s': %w", filePath, err)
		}

		return true, errs.Wrap(err)
	})

	return parsedSignalsCount, errs.Wrap(err)
}

type splittingConsumer struct {
	fileName    string
	dirPath     string
	getConsumer SignalConsumerSource
	consumers   map[string]map[string]ClosableSignalConsumer
}

func (this *splittingConsumer) getConsumerForSignal(s signal.Signal) (SignalConsumer, error) {
	devices, ok := this.consumers[s.Device]
	if !ok {
		consumer, err := this.getConsumer(filepath.Join(this.dirPath, s.Device, fmt.Sprintf("%s_%s", s.Brand, this.fileName)))
		if err != nil {
			return nil, errs.Wrap(err)
		}

		this.consumers[s.Device] = map[string]ClosableSignalConsumer{s.Brand: consumer}
		return consumer, nil
	}

	consumer, ok := devices[s.Brand]
	if !ok {
		consumer, err := this.getConsumer(filepath.Join(this.dirPath, s.Device, fmt.Sprintf("%s_%s", s.Brand, this.fileName)))
		if err != nil {
			return nil, errs.Wrap(err)
		}
		devices[s.Brand] = consumer
		return consumer, nil
	}

	return consumer, nil
}

func (this *splittingConsumer) Consume(s signal.Signal) error {
	consumer, err := this.getConsumerForSignal(s)
	if err == nil {
		err = consumer.Consume(s)
	}
	return errs.Wrap(err)
}

func (this *splittingConsumer) Close() error {
	allErrors := make([]error, 0)
	for _, m := range this.consumers {
		for _, consumer := range m {
			e := consumer.Close()
			if e != nil {
				allErrors = append(allErrors, e)
			}
		}
	}
	if len(allErrors) != 0 {
		return errs.Join(allErrors...)
	}
	return nil
}

func ParseCsvFile(filePath string, options Options, consumer SignalConsumer) (int, error) {
	stream, err := fs.OpenReadOnlyFile(filePath)
	if err != nil {
		return 0, errs.Wrap(err)
	}

	defer func() {
		_ = stream.Close()
	}()

	cfg := parseCfg{
		source: filePath,
	}

	c, err := ParseCsvStream(cfg, stream, consumer)
	return c, errs.Wrap(err)
}
