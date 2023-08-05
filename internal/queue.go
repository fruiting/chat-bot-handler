package internal

import "github.com/IBM/sarama"

//go:generate mockgen -source=queue.go -destination=./queue_mock.go -package=internal

type QueueProducer interface {
	Push(message *sarama.ProducerMessage) error
}
