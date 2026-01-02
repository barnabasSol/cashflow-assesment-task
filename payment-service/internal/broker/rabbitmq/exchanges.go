package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewPaymentExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"payment_exchange", // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)

	return err
}
