package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"irptools/signals/signal"
	"irptools/utils/errs"
	"irptools/utils/fs"
)

func EnumSignals(ctx context.Context, rootPath string, getConsumer SignalsToFileConsumerSourceFn) error {
	err := fs.EnumFilePathsWithExt(rootPath, ".json", func(filePath string) (res bool, err error) {

		f, err := fs.OpenReadOnlyFile(filePath)
		if err != nil {
			return false, errs.Wrap(err)
		}
		defer func() {
			err = errs.Join(err, f.Close())
		}()

		consumer, err := getConsumer(filePath)
		if err != nil {
			return false, errs.Wrap(err)
		}
		defer func() {
			err = errs.Join(err, consumer.Close())
		}()

		err = EnumStreamSignals(ctx, f, consumer)
		if err != nil {
			return false, errs.Wrap(err)
		}

		return true, nil
	})

	return errs.Wrap(err)
}

func EnumStreamSignals(ctx context.Context, stream io.Reader, consumer SignalConsumer) error {
	decoder := json.NewDecoder(stream)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			break
		}

		s := signal.Signal{}
		err := decoder.Decode(&s)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errs.Wrap(err)
		}

		err = consumer.Consume(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}
