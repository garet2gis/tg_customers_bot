package bot_router

import (
	"context"
	cs "github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotRouter) handleCreateServiceBranch(message *tgbotapi.Message, curState cs.State) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Error")

	newState, err := b.serviceHandler.CreateServiceHandler(context.TODO(), message, curState)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	msg = newState.Message
	if newState.Step < 0 {
		if err = b.chatState.Delete(message.From.ID, cs.ChatStateBucket); err != nil {
			return tgbotapi.MessageConfig{}, err
		}
	} else {
		if err = b.chatState.Update(message.From.ID, createStepString(CreateServiceBranch, newState.Step), cs.ChatStateBucket); err != nil {
			return tgbotapi.MessageConfig{}, err
		}
	}

	return msg, nil
}
