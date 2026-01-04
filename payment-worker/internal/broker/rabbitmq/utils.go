package broker

import amqp "github.com/rabbitmq/amqp091-go"

func getRetryCount(msg *amqp.Delivery) int {
	if xDeath, ok := msg.Headers["x-death"].([]interface{}); ok && len(xDeath) > 0 {
		if deathStats, ok := xDeath[0].(amqp.Table); ok {
			if count, ok := deathStats["count"].(int64); ok {
				return int(count)
			}
		}
	}
	return 0
}
