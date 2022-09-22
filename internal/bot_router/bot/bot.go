package bot_router

import (
	"github.com/garet2gis/tg_customers_bot/internal/bot_router"
	"github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"github.com/garet2gis/tg_customers_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotRouter struct {
	bot            telegram.Bot
	logger         *logging.Logger
	chatState      chat_repository.ChatRepository
	serviceHandler bot_router.ServiceHandler
}

func NewBot(bot telegram.Bot, logger *logging.Logger, chatState chat_repository.ChatRepository, serviceHandler bot_router.ServiceHandler) bot_router.BotRouter {
	return &BotRouter{
		bot:            bot,
		logger:         logger,
		serviceHandler: serviceHandler,
		chatState:      chatState,
	}
}

func (b *BotRouter) handleUpdates(updates tgbotapi.UpdatesChannel) {
	var err error
	for update := range updates {

		if update.Message != nil {
			b.logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.IsCommand() {
				err = b.handleCommand(update.Message)
			} else {
				err = b.handleMessage(update.Message)
			}
			if err != nil {
				// TODO: handle error
			}
		} else if update.CallbackQuery != nil {
			err = b.handleCallBackQuery(update.CallbackQuery)
			if err != nil {
				// TODO: handle error
			}
		}
	}

}

func (b *BotRouter) initUpdatesChanel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *BotRouter) Start() {
	updates := b.initUpdatesChanel()
	b.handleUpdates(updates)
}
