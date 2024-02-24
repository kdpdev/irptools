package utils

import "irptools/signals/signal"

func NewPreprocessingSignalConsumer(origin ClosableSignalConsumer, process func(signal.Signal) (signal.Signal, error)) *PreprocessingSignalConsumer {
	return &PreprocessingSignalConsumer{
		origin:  origin,
		process: process,
	}
}

type PreprocessingSignalConsumer struct {
	origin  ClosableSignalConsumer
	process func(signal.Signal) (signal.Signal, error)
}

func (this *PreprocessingSignalConsumer) Consume(s signal.Signal) error {
	s, err := this.process(s)
	if err != nil {
		return err
	}
	return this.origin.Consume(s)
}

func (this *PreprocessingSignalConsumer) Close() error {
	return this.origin.Close()
}
