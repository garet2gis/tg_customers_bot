package bot_router

import (
	"context"
	"errors"
	"fmt"
	cs "github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	chat_storage "github.com/garet2gis/tg_customers_bot/internal/chat_repository/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

const (
	startCommand         = "start"
	createServiceMessage = "Добавить услугу"
)

const CreateServiceBranch = "create_service"

const (
	createService1 = CreateServiceBranch + ".1"
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(createServiceMessage),
	),
)

func (b *BotRouter) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Тестовое сообщение")
	c := message.Command()
	switch c {
	case startCommand:
		msg = tgbotapi.NewMessage(message.Chat.ID, "Возможности:")
		msg.ReplyMarkup = mainKeyboard
	default:
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("команды %s нет", c))
	}
	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("message handler error: %v", err)
	}
	return nil
}

func (b *BotRouter) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "test")

	// проверка на существование истории переписки
	curState, err := b.chatState.Get(message.From.ID, cs.ChatStateBucket)
	if err != nil {
		if errors.Is(err, chat_storage.NoChatStateFound) {
			curState = nil
		} else {
			return err
		}
	}

	switch {
	case message.Text == createServiceMessage:
		msg = tgbotapi.NewMessage(message.Chat.ID, "Введите название")
		if err = b.chatState.Update(message.From.ID, createService1, cs.ChatStateBucket); err != nil {
			return err
		}
	case curState == nil:
		msg = tgbotapi.NewMessage(message.Chat.ID, "Выберите действие")
	case curState.Branch == CreateServiceBranch:
		msg, err = b.handleCreateServiceBranch(message, *curState)
		if err != nil {
			return err
		}
	}

	_, err = b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("message handler error: %v", err)
	}
	return nil
}

func (b *BotRouter) handleCallBackQuery(update *tgbotapi.CallbackQuery) error {
	callback := tgbotapi.NewCallback(update.ID, "данные")
	if _, err := b.bot.Request(callback); err != nil {
		return fmt.Errorf("message handler error (request): %v", err)
	}
	return nil
}

func (b *BotRouter) handleCreateServiceBranch(message *tgbotapi.Message, curState cs.State) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Error")

	newState, err := b.serviceHandler.CreateServiceHandler(context.TODO(), message, curState)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	msg = newState.Message
	if newState.Step < 0 {
		if err = b.chatState.Delete(message.From.ID, cs.ChatStateBucket); err != nil {
			return tgbotapi.MessageConfig{}, err
		}
	} else {
		if err = b.chatState.Update(message.From.ID, createStepString(CreateServiceBranch, newState.Step), cs.ChatStateBucket); err != nil {
			return tgbotapi.MessageConfig{}, err
		}
	}

	return msg, nil
}

func createStepString(branch string, step int) string {
	return branch + "." + strconv.Itoa(step)
}
