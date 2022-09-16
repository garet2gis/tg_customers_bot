package main

import (
	"tg_customers_bot/internal/config"
	"tg_customers_bot/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Debug("Start application")
	cfg := config.GetConfig()
	logger.Debugf("Config: %+v", cfg)
}
