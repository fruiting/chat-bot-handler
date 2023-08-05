package telegram

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type chatBotHandlerSuite struct {
	suite.Suite

	logs    *observer.ObservedLogs
	testErr error
	apiKey  string

	client *MockhttpClient

	handler *ChatBotHandler
}

func TestChatBotHandlerSuite(t *testing.T) {
	suite.Run(t, &chatBotHandlerSuite{})
}

func (s *chatBotHandlerSuite) SetupTest() {
	core, logs := observer.New(zap.InfoLevel)
	s.logs = logs
	s.testErr = errors.New("test err")
	s.apiKey = "test"

	ctrl := gomock.NewController(s.T())
	s.client = NewMockhttpClient(ctrl)

	s.handler = NewChatBotHandler(s.client, s.apiKey, zap.New(core))
}

func (s *chatBotHandlerSuite) TestSendMessagePostErr() {
	s.client.
		EXPECT().
		Post(
			fmt.Sprintf("%s%s/%s", tgUrl, s.apiKey, sendMessageUrl),
			[]byte("{\"chat_id\":1234567890,\"text\":\"test\"}"),
		).
		Return(s.testErr)

	err := s.handler.SendMessage(1234567890, "test")

	s.Equal(fmt.Errorf("can't send message: %w", s.testErr), err)
}

func (s *chatBotHandlerSuite) TestSendMessagePosOk() {
	s.client.
		EXPECT().
		Post(
			fmt.Sprintf("%s%s/%s", tgUrl, s.apiKey, sendMessageUrl),
			[]byte("{\"chat_id\":1234567890,\"text\":\"test\"}"),
		).
		Return(nil)

	err := s.handler.SendMessage(1234567890, "test")

	s.Nil(err)
}
