package main

import (
	"context"
	"flag"
	"github.com/garet2gis/tg_customers_bot/internal/config"
	"github.com/garet2gis/tg_customers_bot/internal/service"
	"github.com/garet2gis/tg_customers_bot/internal/service/db"
	"github.com/garet2gis/tg_customers_bot/pkg/client/postgresql"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"github.com/garet2gis/tg_customers_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Debug("Start application")

	cfg := config.GetConfig()
	logger.Debugf("Config: %+v", cfg)

	token := mustToken(logger)

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), logger, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	serRep := db.NewRepository(postgreSQLClient, logger)

	// create
	newPS := &service.PaidService{
		ID:           "",
		Name:         "TEST",
		BaseDuration: time.Hour,
	}
	err = serRep.Create(context.TODO(), newPS)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	logger.Infof("One: %v", *newPS)

	// find all
	all, err := serRep.FindAll(context.TODO())
	if err != nil {
		logger.Fatalf("%v", err)
	}
	for _, ser := range all {
		logger.Infof("%v", ser)
	}

	// find one
	one, err := serRep.FindOne(context.TODO(), newPS.ID)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	logger.Infof("One: %v", *one)

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
