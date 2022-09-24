package paid_service

import (
	"context"
	"errors"
	"github.com/garet2gis/tg_customers_bot/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
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

func (bs *BotService) CreateServiceStep1(message *tgbotapi.Message, branch string) (string, error) {
	var service = &CreatePaidServiceDTO{Name: message.Text}
	err := bs.serviceTemporaryRepository.Update(message.From.ID, service, branch)
	if err != nil {
		return "", err
	}
	return "Введите длительность в формате: []ч[]м", nil
}

func parseRussianDurationString(russianDuration string) (dur time.Duration, err error) {
	replacer := strings.NewReplacer(" ", "", "д", "d", "ч", "h", "м", "m", "с", "s")
	srtDur := replacer.Replace(russianDuration)
	dur, err = time.ParseDuration(srtDur)
	if err != nil {
		return dur, err
	}
	return dur, nil
}

var WrongDurationFormat = errors.New("wrong duration format")

func (bs *BotService) CreateServiceStep2(ctx context.Context, message *tgbotapi.Message, branch string) (string, error) {
	serviceDTO, err := bs.serviceTemporaryRepository.Get(message.Chat.ID, branch)
	if err != nil {
		return "", err
	}

	duration, err := parseRussianDurationString(message.Text)
	if err != nil {
		return "", WrongDurationFormat
	}

	service := PaidService{
		ID:           "",
		Name:         serviceDTO.Name,
		BaseDuration: duration,
	}
	_, err = bs.service.Create(ctx, &service)
	if err != nil {
		return "Ошибка создания", nil
	}

	err = bs.serviceTemporaryRepository.Delete(message.Chat.ID, branch)
	if err != nil {
		return "", err
	}
	return "Услуга успешно создана!", nil
}

func (bs *BotService) ShowServices(ctx context.Context) (string, error) {
	var msg string

	all, err := bs.service.FindAll(ctx)
	if err != nil {
		return "", err
	}
	for i, s := range all {
		msg += strconv.Itoa(i+1) + ") " + s.String() + "\n"
	}

	return msg, nil
}

func (bs *BotService) DeleteService(ctx context.Context, message *tgbotapi.Message) string {
	var msg = "Запись была успешно удалена"

	id, err := strconv.Atoi(message.Text)

	if err != nil {
		return "Некорректный номер услуги"
	}

	err = bs.service.Delete(ctx, id-1)
	if err != nil {
		if errors.Is(NoRowsDeleted, err) {
			return "Ничего не было удалено"
		}
		if errors.Is(CanNotDeleteRowForeignKey, err) {
			return "Не могу удалить, с данной услугой есть записи"
		}
		return "Ошибка удаления"
	}

	return msg
}
