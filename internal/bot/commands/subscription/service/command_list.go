package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/bot/path"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

func (c *serviceCommander) List(ctx context.Context, inputMessage *tgbotapi.Message, offset uint64, limit uint64) {
	services, err := c.service.List(ctx, offset, limit)
	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.List: failed getting list of services", "err", err)
		return
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, createListText(services.Services))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			createListButtons(
				ctx,
				services.IsHasPreviousPage(),
				services.IsHasNextPage(),
				offset,
				limit,
			)...))

	if _, err := c.bot.Send(msg); err != nil {
		logger.ErrorKV(ctx, "serviceCommander.List: error sending reply message to chat", "err", err)
	}
}

func createListText(services []*subscription.Service) string {
	var sb strings.Builder

	sb.WriteString("Here all the services: \n\n")

	for _, service := range services {
		sb.WriteString(fmt.Sprintf("%v\n", service))
	}

	return sb.String()
}

func createListButtons(ctx context.Context, isHasPreviousPage bool, isHasNextPage bool, offset uint64, limit uint64) []tgbotapi.InlineKeyboardButton {

	buttons := make([]tgbotapi.InlineKeyboardButton, 0)

	serializedData, err := jsonMarshalCallbackListData(
		ctx,
		CallbackListData{
			Offset: offset,
			Limit:  limit,
		},
	)

	if err == nil {
		callbackPath := path.CallbackPath{
			Domain:       "subscription",
			Subdomain:    "service",
			CallbackName: "list",
			CallbackData: string(serializedData),
		}

		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Reload", callbackPath.String()))
	}

	if isHasPreviousPage {
		serializedData, err := jsonMarshalCallbackListData(
			ctx,
			CallbackListData{
				Offset: offset - limit,
				Limit:  limit,
			},
		)

		if err == nil {
			callbackPath := path.CallbackPath{
				Domain:       "subscription",
				Subdomain:    "service",
				CallbackName: "list",
				CallbackData: string(serializedData),
			}

			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("<< Prev page", callbackPath.String()))
		}
	}

	if isHasNextPage {
		serializedData, err := jsonMarshalCallbackListData(
			ctx,
			CallbackListData{
				Offset: offset + limit,
				Limit:  limit,
			},
		)

		if err == nil {
			callbackPath := path.CallbackPath{
				Domain:       "subscription",
				Subdomain:    "service",
				CallbackName: "list",
				CallbackData: string(serializedData),
			}

			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Next page >>", callbackPath.String()))
		}
	}

	return buttons
}

func jsonMarshalCallbackListData(ctx context.Context, data CallbackListData) ([]byte, error) {
	serializedData, err := json.Marshal(data)

	if err != nil {
		logger.ErrorKV(ctx, "jsonMarshalCallbackListData: failed serializing json data", "err", err, "data", data)
		return nil, err
	}

	return serializedData, nil
}
