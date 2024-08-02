package exchange

import "github.com/streadway/amqp"

const CasperEventName = "casper-events"

type Consumer interface {
	Consume(func([]byte)) error
}

type DefaultConsumer struct {
	ch        *amqp.Channel
	queueName string
}

func NewConsumer(queueName string, ch *amqp.Channel) *DefaultConsumer {
	return &DefaultConsumer{ch, queueName}
}

func (c *DefaultConsumer) Consume(handlerFunc func([]byte)) error {
	message, err := c.ch.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range message {
			handlerFunc(d.Body)
		}
	}()

	<-forever

	return nil
}
