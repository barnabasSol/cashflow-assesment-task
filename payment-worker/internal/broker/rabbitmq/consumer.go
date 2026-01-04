package broker

import (
	"payment-worker/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) SubscribeNewPayment(que_name, binding_key string) (<-chan amqp.Delivery, error) {
	retryQueName := que_name + "_retry"
	_, err := r.ch.QueueDeclare(
		retryQueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    PaymentExchangeName,
			"x-dead-letter-routing-key": PaymentExchangeRoutingKey,
			"x-message-ttl":             config.GetEnvInt32("DLX_GAP", 10_000),
		},
	)
	if err != nil {
		return nil, err
	}

	mainArgs := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": retryQueName,
	}

	q, err := r.ch.QueueDeclare(
		que_name,
		true,
		false,
		false,
		false,
		mainArgs,
	)
	if err != nil {
		return nil, err
	}

	err = r.ch.QueueBind(
		q.Name,
		PaymentExchangeRoutingKey,
		PaymentExchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return r.ch.Consume(q.Name, "payment-consumer", false, false, false, false, nil)
}
