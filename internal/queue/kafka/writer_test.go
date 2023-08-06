package kafka

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"fruiting/chat-bot-handler/internal/mock"
	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type writerSuite struct {
	suite.Suite

	topic       string
	msg         []byte
	now         time.Time
	producerMsg *sarama.ProducerMessage
	testErr     error

	producer *mock.MockSyncProducer

	writer *Writer
}

func TestWriterSuite(t *testing.T) {
	suite.Run(t, &writerSuite{})
}

func (s *writerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.topic = "test.topic.v1"
	s.msg = []byte("123")
	s.now = time.Now()
	s.producerMsg = &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.StringEncoder(s.msg),
		Timestamp: s.now,
	}
	s.testErr = errors.New("test err")

	s.producer = mock.NewMockSyncProducer(ctrl)

	s.writer = NewWriter(s.producer)
}

func (s *writerSuite) TestPushErr() {
	s.producer.EXPECT().SendMessage(s.producerMsg).Return(int32(0), int64(0), s.testErr)

	err := s.writer.Push(s.topic, s.msg, s.now)

	s.Equal(fmt.Errorf("can't push message to kafka topic: %w", s.testErr), err)
}

func (s *writerSuite) TestPushOk() {
	s.producer.EXPECT().SendMessage(s.producerMsg).Return(int32(1), int64(2), nil)

	err := s.writer.Push(s.topic, s.msg, s.now)

	s.Nil(err)
}
