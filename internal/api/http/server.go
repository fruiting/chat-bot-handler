package http

//go:generate mockgen -source=server.go -destination=./server_mock.go -package=http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"

	"go.uber.org/zap"
)

type chatBotProcessor interface {
	Process(body []byte) error
}

const pprofUrlPrefix = "/debug/pprof/"

type Server struct {
	listen           string
	chatBotProcessor chatBotProcessor
	enablePprof      bool
	logger           *zap.Logger
}

func NewServer(listen string, chatBotProcessor chatBotProcessor, enablePprof bool, logger *zap.Logger) *Server {
	return &Server{
		listen:           listen,
		chatBotProcessor: chatBotProcessor,
		enablePprof:      enablePprof,
		logger:           logger,
	}
}

func (s *Server) ListenAndServe() error {
	r := http.NewServeMux()

	r.HandleFunc("/", s.handleChatBot)
	r.HandleFunc("/api/v1/ping/", s.ping)

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
		s.logger.Error("can't ready body", zap.Error(err))

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
