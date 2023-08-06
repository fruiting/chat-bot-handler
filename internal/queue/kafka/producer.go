package kafka

//go:generate mockgen -source=writer.go -destination=./writer_mock.go -package=kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(producer sarama.SyncProducer) *Producer {
	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Push(topic string, msg []byte, now time.Time) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(msg),
		Timestamp: now,
	})
	if err != nil {
		return fmt.Errorf("can't push message to kafka topic: %w", err)
	}

	return nil
}
