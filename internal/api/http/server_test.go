package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"fruiting/chat-bot-handler/internal"
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

	chatBotHandler   *internal.MockChatBotHandler
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
	s.chatBotHandler = internal.NewMockChatBotHandler(ctrl)
	s.chatBotProcessor = NewMockchatBotProcessor(ctrl)

	s.server = NewServer(":8080", s.chatBotHandler, s.chatBotProcessor, true, zap.New(core))
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

func (s *serverSuite) TestPing() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ping", nil)

	s.server.ping(s.writer, req)

	body, err := io.ReadAll(s.writer.Body)
	s.Equal(http.StatusOK, s.writer.Code)
	s.Equal("PONG", string(body))
	s.Nil(err)
}

func (s *serverSuite) TestSendMessageUnmarshalErr() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/send-message/", nil)

	s.server.sendMessage(s.writer, req)

	body, err := io.ReadAll(s.writer.Body)
	s.Equal(http.StatusInternalServerError, s.writer.Code)
	s.Equal("", string(body))
	s.Nil(err)
}

func (s *serverSuite) TestSendMessageChatIdIsRequiredErr() {
	type rawReq struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}
	jsonReq, err := json.Marshal(rawReq{
		Text: "text",
	})
	s.Nil(err)

	buf := bytes.Buffer{}
	buf.Write(jsonReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/send-message/", &buf)

	s.server.sendMessage(s.writer, req)

	body, err := io.ReadAll(s.writer.Body)
	s.Equal(http.StatusBadRequest, s.writer.Code)
	s.Equal("chat_id is required", string(body))
	s.Nil(err)
}

func (s *serverSuite) TestSendMessageTextIsRequiredErr() {
	type rawReq struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}
	jsonReq, err := json.Marshal(rawReq{
		ChatId: 1,
	})
	s.Nil(err)

	buf := bytes.Buffer{}
	buf.Write(jsonReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/send-message/", &buf)

	s.server.sendMessage(s.writer, req)

	body, err := io.ReadAll(s.writer.Body)
	s.Equal(http.StatusBadRequest, s.writer.Code)
	s.Equal("text is required", string(body))
	s.Nil(err)
}

func (s *serverSuite) TestSendMessageErr() {
	chatId := int64(1)
	text := "123"

	type rawReq struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}
	jsonReq, err := json.Marshal(rawReq{
		ChatId: chatId,
		Text:   text,
	})
	s.Nil(err)

	buf := bytes.Buffer{}
	buf.Write(jsonReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/send-message/", &buf)
	s.chatBotHandler.EXPECT().SendMessage(internal.ChatId(chatId), internal.Text(text)).Return(s.testErr)

	s.server.sendMessage(s.writer, req)

	body, err := io.ReadAll(s.writer.Body)
	s.Equal(http.StatusInternalServerError, s.writer.Code)
	s.Equal("", string(body))
	s.Nil(err)
}

func (s *serverSuite) TestSendMessageOk() {
	chatId := int64(1)
	text := "123"

	type rawReq struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}
	jsonReq, err := json.Marshal(rawReq{
		ChatId: chatId,
		Text:   text,
	})
	s.Nil(err)

	buf := bytes.Buffer{}
	buf.Write(jsonReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/send-message/", &buf)
	s.chatBotHandler.EXPECT().SendMessage(internal.ChatId(chatId), internal.Text(text)).Return(nil)

	s.server.sendMessage(s.writer, req)

	body, err := io.ReadAll(s.writer.Body)
	s.Equal(http.StatusOK, s.writer.Code)
	s.Equal("", string(body))
	s.Nil(err)
}
