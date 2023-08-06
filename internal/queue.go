package internal

import (
	"time"
)

//go:generate mockgen -source=queue.go -destination=./queue_mock.go -package=internal

type QueueProducer interface {
	Push(topic string, msg []byte, now time.Time) error
}
