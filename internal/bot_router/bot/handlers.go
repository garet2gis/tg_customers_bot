package bot_router

import (
	"errors"
	"fmt"
	cs "github.com/garet2gis/tg_customers_bot/internal/chat_repository"
	chat_storage "github.com/garet2gis/tg_customers_bot/internal/chat_repository/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func (b *BotRouter) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "bot_router_command")
	c := message.Command()
	switch c {
	case startCommand:
		msg = tgbotapi.NewMessage(message.Chat.ID, "Возможности:")
		msg.ReplyMarkup = mainKeyboard
	case stopCommand:
		if err := b.chatState.Delete(message.Chat.ID, cs.ChatStateBucket); err != nil {
			return err
		}
		msg = tgbotapi.NewMessage(message.Chat.ID, "Прервано")
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
	msg := tgbotapi.NewMessage(message.Chat.ID, "bot_router_message")

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
	// проверка на новые действия
	case message.Text == createServiceMessage:
		msg = tgbotapi.NewMessage(message.Chat.ID, "Введите название")
		if err = b.chatState.Update(message.From.ID, createService1, cs.ChatStateBucket); err != nil {
			return err
		}
	case message.Text == updateServiceMessage:
		msg, err = b.handleShowServicesBranch(message)
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Список услуг:\n%sВведите номер услуги, которую желаете изменить", msg.Text))
		if err = b.chatState.Update(message.From.ID, updateService1, cs.ChatStateBucket); err != nil {
			return err
		}
	case message.Text == showServiceMessage:
		msg, err = b.handleShowServicesBranch(message)
		if err != nil {
			return err
		}
	case message.Text == deleteServiceMessage:
		msg, err = b.handleShowServicesBranch(message)
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Список услуг:\n%sВведите номер услуги, которую желаете удалить", msg.Text))
		if err = b.chatState.Update(message.From.ID, deleteService1, cs.ChatStateBucket); err != nil {
			return err
		}
	// если пользователь ввел незначащий текст
	case curState == nil:
		msg = tgbotapi.NewMessage(message.Chat.ID, "Выберите действие")
	// если пользователь находится в ветке диалога
	case curState.Branch == CreateServiceBranch:
		msg, err = b.handleCreateServiceBranch(message, *curState)
		if err != nil {
			return err
		}
	case curState.Branch == UpdateServiceBranch:
		msg, err = b.handleUpdateServiceBranch(message, *curState)
		if err != nil {
			return err
		}
	case curState.Branch == DeleteServiceBranch:
		msg, err = b.handleDeleteServiceBranch(message)
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

func createStepString(branch string, step int) string {
	return branch + "." + strconv.Itoa(step)
}
