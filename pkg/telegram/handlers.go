package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const startCommand = "start"

func (b *Bot) handleCommand(command *tgbotapi.Message) error {
	var msg tgbotapi.MessageConfig
	c := command.Command()
	switch c {
	case startCommand:
		msg = tgbotapi.NewMessage(command.Chat.ID, fmt.Sprintf("это команда: %s", c))
	default:
		msg = tgbotapi.NewMessage(command.Chat.ID, fmt.Sprintf("команды %s нет", c))
	}
	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("command handler error: %v", err)
	}
	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("message handler error: %v", err)
	}
	return nil
}
