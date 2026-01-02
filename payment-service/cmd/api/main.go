package main

import (
	"context"
	broker "payment-service/internal/broker/rabbitmq"
	"payment-service/internal/db"
	"payment-service/internal/server"
)

func main() {
	pg := db.InitPostgres(context.Background())
	defer pg.ConnPool.Close()
	rmq, err := broker.InitRabbitMQ()
	if err != nil {
		panic(err)
	}
	defer rmq.Close()
	server := server.New(":1000", pg, rmq)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
