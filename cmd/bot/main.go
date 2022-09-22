package main

import (
	"context"
	"flag"
	"github.com/garet2gis/tg_customers_bot/internal/bot_router"
	br "github.com/garet2gis/tg_customers_bot/internal/bot_router/bot"
	cr "github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	chatRepository "github.com/garet2gis/tg_customers_bot/internal/chat_repository/db"
	"github.com/garet2gis/tg_customers_bot/internal/config"
	"github.com/garet2gis/tg_customers_bot/internal/paid_service"
	"github.com/garet2gis/tg_customers_bot/internal/paid_service/db"
	"github.com/garet2gis/tg_customers_bot/pkg/client/boltdb"
	"github.com/garet2gis/tg_customers_bot/pkg/client/postgresql"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	"github.com/garet2gis/tg_customers_bot/pkg/telegram"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Debug("Start application")

	cfg := config.GetConfig()
	logger.Debugf("Config: %+v", cfg)

	token := mustToken(logger)

	postgresRepository, err := initPostgresRepository(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	boltDBClient, err := boltdb.NewKeyValueClient("bot.db", []string{cr.ChatStateBucket, br.CreateServiceBranch}, logger)
	if err != nil {
		logger.Fatal(err)
	}

	serviceHandler := initServiceHandler(boltDBClient, *postgresRepository, logger)

	bot := initBot(
		token,
		boltDBClient,
		logger,
		serviceHandler,
	)

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

func initPostgresRepository(logger *logging.Logger, cfg *config.Config) (*paid_service.Repository, error) {
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), logger, cfg.Storage)
	if err != nil {
		return nil, err
	}
	repo := db.NewRepository(postgreSQLClient, logger)
	return &repo, nil
}

func initServiceHandler(client boltdb.KeyValueClient, postgresRepository paid_service.Repository, logger *logging.Logger) bot_router.ServiceHandler {
	serviceRepository := db.NewServiceTemporaryRepository(client, logger)
	return paid_service.NewHandler(serviceRepository, postgresRepository, logger)
}

func initBot(token string, client boltdb.KeyValueClient, logger *logging.Logger, serviceHandler bot_router.ServiceHandler) bot_router.BotRouter {
	botAPI, err := telegram.NewBot(token)
	if err != nil {
		logger.Fatal(err)
	}
	chatRepo := chatRepository.NewChatRepository(client, logger)
	return br.NewBot(
		botAPI,
		logger,
		chatRepo,
		serviceHandler,
	)
}
