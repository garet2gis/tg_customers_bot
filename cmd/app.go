package main

import (
	"flag"
	"tg_customers_bot/internal/config"
	"tg_customers_bot/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Debug("Start application")
	cfg := config.GetConfig()
	logger.Debugf("Config: %+v", cfg)
	token := mustToken(logger)
	logger.Debugf("Token: %s", token)
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
