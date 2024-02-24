package utils

import "irptools/signals/signal"

func NewFilteringSignalConsumer(origin ClosableSignalConsumer, filter func(signal.Signal) (bool, error)) *FilteringSignalConsumer {
	return &FilteringSignalConsumer{
		origin: origin,
		filter: filter,
	}
}

type FilteringSignalConsumer struct {
	origin ClosableSignalConsumer
	filter func(signal.Signal) (bool, error)
}

func (this *FilteringSignalConsumer) Consume(s signal.Signal) error {
	pass, err := this.filter(s)
	if err != nil {
		return err
	}
	if pass {
		return this.origin.Consume(s)
	}
	return nil
}

func (this *FilteringSignalConsumer) Close() error {
	return this.origin.Close()
}
