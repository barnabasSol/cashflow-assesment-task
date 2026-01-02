package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	broker "payment-service/internal/broker/rabbitmq"
	"payment-service/internal/db"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	addr string
	echo *echo.Echo
	db   *db.Postgres
	rmq  *broker.RabbitMQ
}

func New(
	addr string,
	db *db.Postgres,
	rmq *broker.RabbitMQ,
) *Server {
	return &Server{
		addr: addr,
		echo: echo.New(),
		db:   db,
		rmq:  rmq,
	}
}

func (s *Server) Run() error {
	s.echo.Use(middleware.RequestLogger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.ContextTimeout(10 * time.Second))

	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "payment service is ok")
	})

	srv := &http.Server{
		Addr:         s.addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := s.echo.StartServer(srv); err != nil {
			log.Fatalf(
				"failed to start the streaming server %v",
				err,
			)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
	)

	s.bootstrap()

	for _, route := range s.echo.Routes() {
		fmt.Printf(
			"%s \t %s\n",
			route.Method,
			route.Path,
		)
	}

	<-quit
	log.Println("payment service is shutting down")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	return s.echo.Shutdown(ctx)

}
