package bot_router

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// commands
const (
	startCommand = "start"
)

// message commands
const (
	createServiceMessage = "Добавить услугу"
)

// branches
const (
	CreateServiceBranch = "create_service"
)

// branch + step
const (
	createService1 = CreateServiceBranch + ".1"
)

// keyboards
var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(createServiceMessage),
	),
)
