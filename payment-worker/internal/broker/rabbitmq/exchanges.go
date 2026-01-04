package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	NewPaymentQue = "new_payment_que"
)

const (
	PaymentExchangeName = "payment_exchange"
)

const (
	PaymentExchangeRoutingKey = "payment.created"
)

func NewPaymentExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		PaymentExchangeName, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
}
