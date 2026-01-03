package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PaymentExchangeName = "payment_exchange"
)
const (
	PaymentExchangeRoutingKey = "payment.created"
)

func NewPaymentExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		PaymentExchangeName, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)

	return err
}
