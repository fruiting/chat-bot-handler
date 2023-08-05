package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Writer struct {
	producer sarama.SyncProducer
}

func NewKafkaWriter(producer sarama.SyncProducer) *Writer {
	return &Writer{
		producer: producer,
	}
}

func (w *Writer) Push(message *sarama.ProducerMessage) error {
	_, _, err := w.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("can't push message to kafka topic: %w", err)
	}

	return nil
}
