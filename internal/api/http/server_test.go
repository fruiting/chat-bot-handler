package http

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type serverSuite struct {
	suite.Suite

	logs    *observer.ObservedLogs
	testErr error
	writer  *httptest.ResponseRecorder

	chatBotProcessor *MockchatBotProcessor

	server *Server
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, &serverSuite{})
}

func (s *serverSuite) SetupTest() {
	core, logs := observer.New(zap.InfoLevel)
	s.logs = logs
	s.testErr = errors.New("test err")
	s.writer = httptest.NewRecorder()

	ctrl := gomock.NewController(s.T())
	s.chatBotProcessor = NewMockchatBotProcessor(ctrl)

	s.server = NewServer(":8080", s.chatBotProcessor, true, zap.New(core))
}

func (s *serverSuite) TestHandleChatBotProcessErr() {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	s.chatBotProcessor.EXPECT().Process([]byte("")).Return(s.testErr)

	s.server.handleChatBot(s.writer, req)

	s.Equal(
		1,
		s.logs.FilterMessage("can't process chat bot message").FilterField(zap.Error(s.testErr)).Len(),
	)
}

func (s *serverSuite) TestHandleChatBotOk() {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	s.chatBotProcessor.EXPECT().Process([]byte("")).Return(nil)

	s.server.handleChatBot(s.writer, req)

	s.Equal(
		0,
		s.logs.FilterMessage("can't process chat bot message").FilterField(zap.Error(s.testErr)).Len(),
	)
}

func (s *serverSuite) TestPong() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ping", nil)

	s.server.ping(s.writer, req)

	body, err := io.ReadAll(s.writer.Body)
	s.Equal(http.StatusOK, s.writer.Code)
	s.Equal("PONG", string(body))
	s.Nil(err)
}
