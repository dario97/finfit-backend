package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Read(ctx context.Context, messageChannel chan<- Message, errorChannel chan<- error)
}

type consumer struct {
	reader *kafka.Reader
}

func NewConsumer(reader *kafka.Reader) Consumer {
	return &consumer{reader: reader}
}

func (c consumer) Read(ctx context.Context, messageChannel chan<- Message, errorChannel chan<- error) {
	defer c.reader.Close()

	log.Info("hola chau")
	for {
		log.Info("hola dario")
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Info("hola dario 2")
			errorChannel <- errors.New(fmt.Sprintf("error while reading a message: %v", err))
			continue
		}

		log.Info("hola dario 3")

		//var message Message
		//err = json.Unmarshal(m.Value, &message)
		println(string(m.Value))
		if err != nil {
			errorChannel <- err
		}

		//messageChannel <- message
	}
}
