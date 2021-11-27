package service

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

func (c *serviceCommander) Get(ctx context.Context, inputMessage *tgbotapi.Message) {
	args := inputMessage.CommandArguments()

	sendMessage := func(text string) {
		msg := tgbotapi.NewMessage(
			inputMessage.Chat.ID,
			text,
		)

		if _, err := c.bot.Send(msg); err != nil {
			logger.ErrorKV(ctx, "serviceCommander.Get: error sending reply message to chat", "err", err)
		}
	}

	serviceID, err := strconv.Atoi(args)
	if err != nil {
		logger.DebugKV(ctx, "serviceCommander.Get: received wrong args", "args", args)
		sendMessage(fmt.Sprintf("Wrong arguments: \"%v\"! See: /help__subscription__service", args))
		return
	}

	service, err := c.service.Describe(ctx, uint64(serviceID))
	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.Get: failed getting service", "err", err, "ID", serviceID)
		sendMessage(fmt.Sprintf("Fail to get service with ID %d.", serviceID))
		return
	}

	sendMessage(fmt.Sprintf("%v", service))
}
