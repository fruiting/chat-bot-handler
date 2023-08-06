package http

//go:generate mockgen -source=server.go -destination=./server_mock.go -package=http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"

	"fruiting/chat-bot-handler/internal"
	"go.uber.org/zap"
)

type chatBotProcessor interface {
	Process(body []byte) error
}

const pprofUrlPrefix = "/debug/pprof/"

type Server struct {
	listen           string
	chatBotHandler   internal.ChatBotHandler
	chatBotProcessor chatBotProcessor
	enablePprof      bool
	logger           *zap.Logger
}

func NewServer(
	listen string,
	chatBotHandler internal.ChatBotHandler,
	chatBotProcessor chatBotProcessor,
	enablePprof bool,
	logger *zap.Logger,
) *Server {
	return &Server{
		listen:           listen,
		chatBotHandler:   chatBotHandler,
		chatBotProcessor: chatBotProcessor,
		enablePprof:      enablePprof,
		logger:           logger,
	}
}

func (s *Server) ListenAndServe() error {
	r := http.NewServeMux()

	r.HandleFunc("/", s.handleChatBot)
	r.HandleFunc("/api/v1/ping/", s.ping)
	r.HandleFunc("/api/v1/send-message/", s.sendMessage)

	if s.enablePprof {
		r.HandleFunc(pprofUrlPrefix, pprof.Index)
		r.HandleFunc(fmt.Sprintf("%s/cmdline", pprofUrlPrefix), pprof.Cmdline)
		r.HandleFunc(fmt.Sprintf("%s/profile", pprofUrlPrefix), pprof.Profile)
		r.HandleFunc(fmt.Sprintf("%s/symbol", pprofUrlPrefix), pprof.Symbol)
		r.HandleFunc(fmt.Sprintf("%s/trace", pprofUrlPrefix), pprof.Trace)
	}

	return http.ListenAndServe(s.listen, r)
}

func (s *Server) handleChatBot(_ http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		s.logger.Error("can't read body", zap.Error(err))

		return
	}

	err = s.chatBotProcessor.Process(body)
	if err != nil {
		s.logger.Error("can't process chat bot message", zap.Error(err))
	}
}

func (s *Server) ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PONG"))
}

func (s *Server) sendMessage(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		s.logger.Error("can't read body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	rawReq := struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}{}
	err = json.Unmarshal(body, &rawReq)
	if err != nil {
		s.logger.Error("can't unmarshal body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if rawReq.ChatId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("chat_id is required"))

		return
	}
	if rawReq.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("text is required"))

		return
	}

	err = s.chatBotHandler.SendMessage(internal.ChatId(rawReq.ChatId), internal.Text(rawReq.Text))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
