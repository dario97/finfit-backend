package main

import (
	"context"
	"finfit-backend/internal/application"
	"finfit-backend/internal/infrastructure/kafka"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	kafkago "github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Info("starting application...")
	log.Info("configuring kafka - NUEVO 12")

	chMsg := make(chan kafka.Message)
	chErr := make(chan error)
	kafkaReaderConfig := kafkago.ReaderConfig{
		Brokers:         []string{"172.20.0.1:9092"},
		Topic:           "dbserver1.public.expense",
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
		GroupID:         "go-consumer",
		StartOffset:     kafkago.LastOffset,
	}
	consumer := kafka.NewConsumer(kafkago.NewReader(kafkaReaderConfig))

	go func() {
		log.Info("consumiendo")
		consumer.Read(context.Background(), chMsg, chErr)
		log.Info("consumiendo 2")
	}()

	log.Info("consumiendo 3")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//	for {
	//		select {
	//		case <-quit:
	//			log.Info("HOLA")
	//			goto end
	//		case _ = <-chMsg:
	//			log.Info("HOLA")
	//		case err := <-chErr:
	//			log.Error(err)
	//		}
	//	}
	//end:
	log.Info("application started on port 8080")

	e := echo.New()
	app := application.NewApplication(e)
	defer app.Finish()
	app.LoadDependencyConfiguration()
	err := app.Start()
	if err != nil {
		panic(err)
	}
}
