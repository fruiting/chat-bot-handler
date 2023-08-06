package internal

//go:generate easyjson -output_filename=./chat_bot_command_easyjson.go

type ChatId int64
type Text string
type ChatBotCommand string

const (
	ParseJobsInfoChatBotCommand ChatBotCommand = "/parse_jobs_info"
)

var whiteListCommands = []ChatBotCommand{
	ParseJobsInfoChatBotCommand,
}

func IsInCommandsWhitelist(command ChatBotCommand) bool {
	for _, whiteListCommand := range whiteListCommands {
		if command == whiteListCommand {
			return true
		}
	}

	return false
}

// easyjson:json ChatBotCommandInfo
type ChatBotCommandInfo struct {
	ChatId    ChatId
	Command   ChatBotCommand
	Parser    string
	Positions []string
	Keywords  []string
	IsReady   bool
}

func (c *ChatBotCommandInfo) validateCommand() error {
	if c.Command != ParseJobsInfoChatBotCommand {
		return InvalidCommandErr
	}

	return nil
}

type chatBotCommands struct {
	commands map[ChatId]*ChatBotCommandInfo
}

func NewChatBotCommands() *chatBotCommands {
	return &chatBotCommands{
		commands: make(map[ChatId]*ChatBotCommandInfo, 0),
	}
}

func (c *chatBotCommands) clearCommandFromMap(chatId ChatId) {
	_, ok := c.commands[chatId]
	if ok {
		delete(c.commands, chatId)
	}
}
