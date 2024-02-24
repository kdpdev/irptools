package utils

import (
	"io"

	"irptools/signals/signal"
)

type SignalConsumer interface {
	Consume(signal signal.Signal) error
}

type ClosableSignalConsumer interface {
	SignalConsumer
	io.Closer
}

type SignalsToFileConsumerSourceFn = func(filePath string) (ClosableSignalConsumer, error)
