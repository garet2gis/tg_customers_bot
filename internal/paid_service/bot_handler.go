package paid_service

import (
	"context"
	"github.com/garet2gis/tg_customers_bot/internal/bot_router"
	cs "github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ServiceHandler struct {
	botService *BotService
	logger     *logging.Logger
}

func NewHandler(serviceTemporaryRepository ServiceTemporaryRepository, repository Repository, logger *logging.Logger) *ServiceHandler {
	botService := NewBotService(serviceTemporaryRepository, repository, logger)
	return &ServiceHandler{
		botService: botService,
		logger:     logger,
	}
}

func (h *ServiceHandler) CreateServiceHandler(ctx context.Context, message *tgbotapi.Message, chatState cs.State) (bot_router.MessageReply, error) {
	res := bot_router.MessageReply{}

	switch chatState.Step {
	case 1:
		err := h.botService.CreateServiceStep1(message, chatState.Branch)
		if err != nil {
			return bot_router.MessageReply{}, err
		}

		res.Message = tgbotapi.NewMessage(message.Chat.ID, "Введите длительность")
		res.Step = chatState.Step + 1
	case 2:
		err := h.botService.CreateServiceStep2(ctx, message, chatState.Branch)
		if err != nil {
			return bot_router.MessageReply{}, err
		}

		res.Message = tgbotapi.NewMessage(message.Chat.ID, "Услуга добавлена!")
		res.Step = -1
	}

	return res, nil
}
func (h *ServiceHandler) ShowServicesHandler(ctx context.Context, message *tgbotapi.Message) (bot_router.MessageReply, error) {
	res := bot_router.MessageReply{}

	services, err := h.botService.ShowServices(ctx)
	if err != nil {
		return bot_router.MessageReply{}, err
	}
	res.Message = tgbotapi.NewMessage(message.Chat.ID, services)

	return res, nil
}
