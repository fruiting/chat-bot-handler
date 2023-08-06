package kafka

//go:generate mockgen -source=writer.go -destination=./writer_mock.go -package=kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Writer struct {
	producer sarama.SyncProducer
}

func NewWriter(producer sarama.SyncProducer) *Writer {
	return &Writer{
		producer: producer,
	}
}

func (w *Writer) Push(topic string, msg []byte, now time.Time) error {
	_, _, err := w.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(msg),
		Timestamp: now,
	})
	if err != nil {
		return fmt.Errorf("can't push message to kafka topic: %w", err)
	}

	return nil
}
