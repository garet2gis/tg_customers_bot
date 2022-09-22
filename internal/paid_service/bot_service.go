package paid_service

import (
	"context"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type BotService struct {
	service                    *Service
	logger                     *logging.Logger
	serviceTemporaryRepository ServiceTemporaryRepository
}

func NewBotService(serviceTemporaryRepository ServiceTemporaryRepository, repository Repository, logger *logging.Logger) *BotService {
	service := NewService(repository, logger)
	return &BotService{
		service:                    service,
		logger:                     logger,
		serviceTemporaryRepository: serviceTemporaryRepository,
	}
}

func (bs *BotService) CreateServiceStep1(message *tgbotapi.Message, branch string) error {
	var service = &CreatePaidServiceDTO{Name: message.Text}
	err := bs.serviceTemporaryRepository.Update(message.From.ID, service, branch)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BotService) CreateServiceStep2(ctx context.Context, message *tgbotapi.Message, branch string) error {
	serviceDTO, err := bs.serviceTemporaryRepository.Get(message.Chat.ID, branch)
	if err != nil {
		return err
	}
	service := PaidService{
		ID:           "",
		Name:         serviceDTO.Name,
		BaseDuration: time.Hour,
	}
	_, err = bs.service.Create(ctx, &service)
	if err != nil {
		return err
	}

	err = bs.serviceTemporaryRepository.Delete(message.Chat.ID, branch)
	if err != nil {
		return err
	}
	return nil
}
