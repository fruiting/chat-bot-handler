package internal

//go:generate mockgen -source=chat_bot_processor.go -destination=./chat_bot_processor_mock.go -package=internal

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mailru/easyjson"
)

const parseJobsTopic string = "job-parser.parse-jobs.v1"

type ChatBotHandler interface {
	FindChatIdAndText(bodyRequest []byte) (ChatId, Text, error)
	SendMessage(chatId ChatId, text Text) error
}

type ChatBotProcessor struct {
	handler         ChatBotHandler
	chatBotCommands *chatBotCommands
	queueProducer   QueueProducer
	mu              *sync.Mutex
}

func NewChatBotProcessor(handler ChatBotHandler, queueProducer QueueProducer) *ChatBotProcessor {
	return &ChatBotProcessor{
		handler:         handler,
		chatBotCommands: NewChatBotCommands(),
		queueProducer:   queueProducer,
		mu:              &sync.Mutex{},
	}
}

func (p *ChatBotProcessor) Process(body []byte) error {
	chatId, text, err := p.handler.FindChatIdAndText(body)
	if err != nil {
		return fmt.Errorf("can't find chat id: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	chatBotCommand, ok := p.chatBotCommands.commands[chatId]
	if chatBotCommand == nil || !ok {
		if !IsInCommandsWhitelist(ChatBotCommand(text)) {
			return nil
		}

		chatBotCommand = &ChatBotCommandInfo{
			ChatId:  chatId,
			Command: ChatBotCommand(text),
		}
		p.chatBotCommands.commands[chatId] = chatBotCommand

		err = p.handler.SendMessage(chatId, "Choose parser")
		if err != nil {
			p.chatBotCommands.clearCommandFromMap(chatId)
			return fmt.Errorf("can't send message to choose parser")
		}

		return nil
	}

	command := strings.Split(string(text), ":")
	if len(command) < 2 {
		p.chatBotCommands.clearCommandFromMap(chatId)
		return InvalidCommandErr
	}

	var msg string
	switch command[0] {
	case "parser":
		chatBotCommand.Parser = command[1]
		msg = "Now type position"
	case "position":
		chatBotCommand.Positions = []string{command[1]} //todo add many
		msg = "Now type keywords"
	case "keywords":
		keywords := strings.Split(command[1], "|")
		chatBotCommand.Keywords = keywords
		chatBotCommand.IsReady = true
		msg = "Processing..."
	default:
		err = p.handleInvalidCommand(chatId)
		if err != nil {
			return fmt.Errorf("can't handle invalid command: %w", err)
		}

		return InvalidCommandErr
	}

	err = chatBotCommand.validateCommand()
	if err != nil {
		err = p.handleInvalidCommand(chatId)
		if err != nil {
			return fmt.Errorf("can't handle invalid command: %w", err)
		}

		return InvalidCommandErr
	}

	err = p.handler.SendMessage(chatId, Text(msg))
	if err != nil {
		p.chatBotCommands.clearCommandFromMap(chatId)
		return fmt.Errorf("can't send message for building command: %w", err)
	}

	if !chatBotCommand.IsReady {
		return nil
	}

	p.chatBotCommands.clearCommandFromMap(chatId)
	if chatBotCommand.Command == ParseJobsInfoChatBotCommand {
		err = p.pushCommandIntoQueue(chatBotCommand)
		if err != nil {
			return fmt.Errorf("can't push command into queue: %w", err)
		}
	}

	return nil
}

func (p *ChatBotProcessor) pushCommandIntoQueue(chatBotCommand *ChatBotCommandInfo) error {
	jsonCommand, err := easyjson.Marshal(chatBotCommand)
	if err != nil {
		return fmt.Errorf("can't marshal chat bot command")
	}

	err = p.queueProducer.Push(parseJobsTopic, jsonCommand, time.Now())
	if err != nil {
		return fmt.Errorf("can't process command: %w", err)
	}

	return nil
}

func (p *ChatBotProcessor) handleInvalidCommand(chatId ChatId) error {
	p.chatBotCommands.clearCommandFromMap(chatId)

	err := p.handler.SendMessage(chatId, Text(InvalidCommandErr.Error()))
	if err != nil {
		return fmt.Errorf("can't send message to choose parser")
	}

	return nil
}
