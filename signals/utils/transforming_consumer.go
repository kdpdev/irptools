package utils

import (
	"irptools/signals/signal"
	"irptools/utils/errs"
)

type TransformSignalFn func(signal.Signal) (signal.Signal, error)

func NewTransformingSignalConsumer(origin ClosableSignalConsumer, trs []TransformSignalFn) ClosableSignalConsumer {
	if len(trs) == 0 {
		return origin
	}

	return &TransformingSignalConsumer{
		origin: origin,
		trs:    trs,
	}
}

type TransformingSignalConsumer struct {
	origin ClosableSignalConsumer
	trs    []TransformSignalFn
}

func (this *TransformingSignalConsumer) Consume(s signal.Signal) error {
	var err error
	for _, tr := range this.trs {
		s, err = tr(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return errs.Wrap(this.origin.Consume(s))
}

func (this *TransformingSignalConsumer) Close() error {
	return this.origin.Close()
}
