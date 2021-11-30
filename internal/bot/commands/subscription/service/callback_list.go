package service

import (
	"context"
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/path"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

// CallbackListData - payload коллбэка списка ссервисов.
type CallbackListData struct {
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

func (c *serviceCommander) CallbackList(ctx context.Context, callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	parsedData := CallbackListData{}
	err := json.Unmarshal([]byte(callbackPath.CallbackData), &parsedData)
	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.CallbackList: failed reading json data", "err", err, "data", callbackPath.CallbackData)
		return
	}

	_, err = c.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.CallbackList: failed sending answer to callback", "err", err)
	}

	offset := parsedData.Offset
	limit := parsedData.Limit

	services, err := c.service.List(ctx, offset, limit)

	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.CallbackList: failed getting list of services", "err", err)
		return
	}

	for services.Len() == 0 && services.IsHasPreviousPage() && offset-limit >= 0 {
		offset -= limit
		services, err = c.service.List(ctx, offset, limit)
		if err != nil {
			logger.ErrorKV(ctx, "serviceCommander.CallbackList: failed getting list of services", "err", err)
			return
		}
	}

	_, err = c.bot.Send(tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, createListText(services.Services)))
	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.CallbackList: failed sending edit text message to chat", "err", err)
		return
	}

	_, err = c.bot.Send(
		tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID,
			tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					createListButtons(
						ctx,
						services.IsHasPreviousPage(),
						services.IsHasNextPage(),
						offset,
						limit,
					)...))))
	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.CallbackList: failed sending edit keyboard message to chat", "err", err)
	}
}
