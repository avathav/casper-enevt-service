package eventupdate

import (
	"event-service/internal/domain/event/aggregate"
	"event-service/internal/exchange"

	"github.com/streadway/amqp"
)

const ExchangeKey = "event-update"
const QueueName = "event-update"

type Producer struct {
	ch *amqp.Channel
}

func NewProducer(ch *amqp.Channel) *Producer {
	return &Producer{ch}
}

func (p *Producer) Publish(body []byte) error {
	return p.ch.Publish(
		exchange.CasperEventName,
		ExchangeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

type Message struct {
	Event *aggregate.Event
}
