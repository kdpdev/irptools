package fz

import (
	"context"
	"io"
	"path/filepath"
	"strings"

	"irptools/signals/signal"
	"irptools/utils/errs"
	"irptools/utils/fs"
	jsonutils "irptools/utils/json"
	"irptools/utils/logs"
)

type Options struct {
	IgnoreAllUnsupportedProtocolsError      bool     `json:"ignoreAllUnsupportedProtocolError"`
	IgnoreSpecificUnsupportedProtocolsError []string `json:"ignoreSpecificUnsupportedProtocolError"`
}

type SignalConsumer interface {
	Consume(signal signal.Signal) error
}

type ClosableSignalConsumer interface {
	SignalConsumer
	io.Closer
}

type SignalConsumerSource = func(filePath string) (ClosableSignalConsumer, error)

type LinesBatch = []string
type FieldsMap = map[string]string

func ParseIrFiles(
	ctx context.Context,
	rootPath string,
	opts interface{},
	getConsumer SignalConsumerSource) (int, error) {

	l := logs.L(ctx)

	options := Options{}
	err := jsonutils.Cast(opts, &options)
	if err != nil {
		return 0, errs.Wrap(err)
	}

	parsedSignalsCount := 0
	err = fs.EnumFilePathsWithExt(rootPath, ".ir", func(filePath string) (res bool, err error) {
		consumer, err := getConsumer(filePath)

		if err != nil {
			return false, errs.Wrap(err)
		}

		defer func() {
			closeErr := consumer.Close()
			if err == nil {
				err = closeErr
			}
		}()

		count, err := ParseIrFile(filePath, options, consumer)
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

func ParseIrFile(filePath string, options Options, consumer SignalConsumer) (int, error) {
	stream, err := fs.OpenReadOnlyFile(filePath)
	if err != nil {
		return 0, errs.Wrap(err)
	}

	defer func() {
		_ = stream.Close()
	}()

	cfg := parseCfg{
		source:                             filePath,
		brand:                              brandFromFilePath(filePath),
		ignoreAllUnsupportedProtocols:      options.IgnoreAllUnsupportedProtocolsError,
		ignoreSpecificUnsupportedProtocols: options.IgnoreSpecificUnsupportedProtocolsError,
	}

	c, err := ParseIrStream(cfg, stream, consumer)
	return c, errs.Wrap(err)
}

func brandFromFilePath(filePath string) string {
	fileName := filepath.Base(filePath)
	separatorPos := strings.IndexAny(fileName, "_-.")
	if separatorPos > 0 {
		return strings.ToLower(fileName[0:separatorPos])
	}
	return fileName
}
