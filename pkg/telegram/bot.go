package telegram

import (
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}

type Bot struct {
	bot    IBot
	logger *logging.Logger
}

func NewBot(bot IBot, logger *logging.Logger) *Bot {
	return &Bot{bot: bot, logger: logger}
}

func (b *Bot) Start() {
	updates := b.initUpdatesChanel()
	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	var err error
	for update := range updates {
		if update.Message == nil {
			continue
		}
		b.logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			err = b.handleCommand(update.Message)
		} else {
			err = b.handleMessage(update.Message)
		}
		if err != nil {
			// TODO: handle error
		}
	}
}

func (b *Bot) initUpdatesChanel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
