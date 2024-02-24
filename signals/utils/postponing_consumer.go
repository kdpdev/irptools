package utils

import "irptools/signals/signal"

func NewPostponingConsumer(createConsumer func() (ClosableSignalConsumer, error)) *PostponingConsumer {
	consumer := &PostponingConsumer{
		createConsumer: createConsumer,
	}
	consumer.consumeFn = consumer.consumeCreate
	return consumer
}

type PostponingConsumer struct {
	consumer       ClosableSignalConsumer
	createConsumer func() (ClosableSignalConsumer, error)
	consumeFn      func(signal.Signal) error
}

func (this *PostponingConsumer) Consume(s signal.Signal) error {
	return this.consumeFn(s)
}

func (this *PostponingConsumer) Close() error {
	if this.consumer != nil {
		return this.consumer.Close()
	}
	return nil
}

func (this *PostponingConsumer) consumeCreate(s signal.Signal) error {
	consumer, err := this.createConsumer()
	if err != nil {
		return err
	}
	this.consumer = consumer
	this.consumeFn = this.consume
	return this.consume(s)
}

func (this *PostponingConsumer) consume(s signal.Signal) error {
	return this.consumer.Consume(s)
}
