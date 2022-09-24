package bot_router

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// commands
const (
	startCommand = "start"
	stopCommand  = "stop"
)

// message commands
const (
	createServiceMessage = "Добавить услугу"
	showServiceMessage   = "Показать услуги"
	deleteServiceMessage = "Удалить услугу"
	updateServiceMessage = "Обновить услугу"
)

// branches
const (
	CreateServiceBranch = "create_service"
	DeleteServiceBranch = "delete_service"
	UpdateServiceBranch = "update_service"
)

// branch + step
const (
	createService1 = CreateServiceBranch + ".1"
	deleteService1 = DeleteServiceBranch + ".1"
	updateService1 = UpdateServiceBranch + ".1"
)

// keyboards
var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(createServiceMessage),
		tgbotapi.NewKeyboardButton(updateServiceMessage),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(showServiceMessage),
		tgbotapi.NewKeyboardButton(deleteServiceMessage),
	),
)
