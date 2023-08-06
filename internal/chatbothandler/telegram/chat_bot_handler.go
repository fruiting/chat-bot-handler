package telegram

//go:generate mockgen -source=chat_bot_handler.go -destination=./chat_bot_handler_mock.go -package=telegram

import (
	"encoding/json"
	"fmt"

	"fruiting/chat-bot-handler/internal"
	"go.uber.org/zap"
)

type httpClient interface {
	Post(url string, body []byte) error
}

const (
	tgUrl          string = "https://api.telegram.org/bot"
	sendMessageUrl string = "sendMessage"
)

type tgRequest struct {
	Message struct {
		Chat struct {
			Id int64 `json:"id"`
		}
		Text string `json:"text"`
	}
}

type ChatBotHandler struct {
	client httpClient
	apiKey string
	logger *zap.Logger
}

func NewChatBotHandler(
	client httpClient,
	apiKey string,
	logger *zap.Logger,
) *ChatBotHandler {
	return &ChatBotHandler{
		client: client,
		apiKey: apiKey,
		logger: logger,
	}
}

func (h *ChatBotHandler) FindChatIdAndText(bodyRequest []byte) (internal.ChatId, internal.Text, error) {
	var request *tgRequest
	err := json.Unmarshal(bodyRequest, &request)
	if err != nil {
		return 0, "", fmt.Errorf("can't unmarshal request: %w", err)
	}

	return internal.ChatId(request.Message.Chat.Id), internal.Text(request.Message.Text), nil
}

func (h *ChatBotHandler) SendMessage(chatId internal.ChatId, text internal.Text) error {
	type sendMessageBody struct {
		ChatId internal.ChatId `json:"chat_id"`
		Text   internal.Text   `json:"text"`
	}

	msg := &sendMessageBody{
		ChatId: chatId,
		Text:   text,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("can't marshal msg: %w", err)
	}

	err = h.client.Post(fmt.Sprintf("%s%s/%s", tgUrl, h.apiKey, sendMessageUrl), msgJson)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}
