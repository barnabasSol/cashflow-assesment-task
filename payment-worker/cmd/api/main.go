package main

import (
	"context"
	"log"
	broker "payment-worker/internal/broker/rabbitmq"
	"payment-worker/internal/db"
)

func main() {
	pg := db.InitPostgres(context.Background())
	defer pg.ConnPool.Close()
	rmq, err := broker.InitRabbitMQ(pg)
	if err != nil {
		panic(err)
	}
	defer rmq.Close()
	new_payment_msgs, err := rmq.SubscribeNewPayment(
		broker.NewPaymentQue,
		"payment.*",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("worker is up")
	rmq.ListenForNewPayment(new_payment_msgs)
}
