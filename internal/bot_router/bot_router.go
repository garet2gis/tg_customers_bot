package bot_router

import (
	"context"
	"github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageReply struct {
	Message tgbotapi.MessageConfig
	Step    int
}

type BotRouter interface {
	Start()
}

type ServiceHandler interface {
	CreateServiceHandler(ctx context.Context, message *tgbotapi.Message, step chat_repository.State) (MessageReply, error)
	UpdateServiceHandler(ctx context.Context, message *tgbotapi.Message, step chat_repository.State) MessageReply
	ShowServicesHandler(ctx context.Context, message *tgbotapi.Message) (MessageReply, error)
	DeleteServiceHandler(ctx context.Context, message *tgbotapi.Message) (MessageReply, error)
}
