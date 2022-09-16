package main

import (
	"flag"
	"github.com/garet2gis/tg_customers_bot/internal/config"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"github.com/garet2gis/tg_customers_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Debug("Start application")

	cfg := config.GetConfig()
	logger.Debugf("Config: %+v", cfg)

	token := mustToken(logger)

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Fatal(err)
	}
	bot := telegram.NewBot(botAPI, logger)
	bot.Start()
}

func mustToken(logger *logging.Logger) string {
	token := flag.String(
		"token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		logger.Fatal("token is not specified")
	}

	return *token
}
