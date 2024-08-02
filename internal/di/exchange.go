package di

import (
	"event-service/internal/config"
	"event-service/internal/exchange"
	"event-service/internal/exchange/event"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func NewAmpqConnection(connectionString string) *amqp.Connection {
	if container.ampqConnection[connectionString] != nil {
		return container.ampqConnection[connectionString]
	}

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		panic(err)
	}

	container.ampqConnection[connectionString] = conn

	return container.ampqConnection[connectionString]
}

func NewAmpqExchange(name string) *amqp.Channel {
	if container.ampqChannels[name] != nil {
		return container.ampqChannels[name]
	}

	channel, chanErr := NewDefaultAmpqConnection().Channel()
	if chanErr != nil {
		log.Panic("Failed to open a channel", chanErr)
	}

	if err := channel.ExchangeDeclare(
		name,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panic("Failed to declare an exchange", err)
	}

	container.ampqChannels[name] = channel

	return container.ampqChannels[name]
}

func NewAmpqQueue(queue, exchangeKey, exchange string) *amqp.Channel {
	channel := NewAmpqExchange(exchange)

	q, err := channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panic("Failed to declare a queue", err)
	}

	if err = channel.QueueBind(
		q.Name,
		exchangeKey,
		exchange,
		false,
		nil,
	); err != nil {
		log.Panic("Failed to bind a queue", err)
	}

	return channel
}

func CloseAmpqConnection(connectionString string) {
	if container.ampqConnection[connectionString] != nil {
		container.ampqConnection[connectionString].Close()
	}
}

func CloseAmpqExchange(name string) {
	if container.ampqChannels[name] != nil {
		container.ampqChannels[name].Close()
	}
}

func CloseAllExchangeConnections() {
	for s, channel := range container.ampqChannels {
		channel.Close()
		delete(container.ampqChannels, s)
	}

	for i, conn := range container.ampqConnection {
		conn.Close()
		delete(container.ampqConnection, i)
	}
}

func NewDefaultAmpqConnection() *amqp.Connection {
	return NewAmpqConnection(config.GetString("AMPQ"))
}

func NewEventUpdateProducer() *eventupdate.Producer {
	return eventupdate.NewProducer(NewAmpqExchange(eventupdate.ExchangeKey))
}

func NewEventUpdateConsumer() *exchange.DefaultConsumer {
	return exchange.NewConsumer(eventupdate.QueueName, NewAmpqQueue(
		eventupdate.QueueName,
		eventupdate.ExchangeKey,
		exchange.CasperEventName,
	))
}
