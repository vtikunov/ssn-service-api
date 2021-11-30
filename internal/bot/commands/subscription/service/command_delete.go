package service

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

func (c *serviceCommander) Delete(ctx context.Context, inputMessage *tgbotapi.Message) {
	args := inputMessage.CommandArguments()

	sendMessage := func(text string) {
		msg := tgbotapi.NewMessage(
			inputMessage.Chat.ID,
			text,
		)

		if _, err := c.bot.Send(msg); err != nil {
			logger.ErrorKV(ctx, "serviceCommander.Delete: error sending reply message to chat", "err", err)
		}
	}

	serviceID, err := strconv.Atoi(args)
	if err != nil {
		logger.DebugKV(ctx, "serviceCommander.Delete: received wrong args", "args", args)
		sendMessage(fmt.Sprintf("Wrong arguments: \"%v\"! See: /help__subscription__service", args))
		return
	}

	if _, err := c.service.Remove(ctx, uint64(serviceID)); err != nil {
		logger.ErrorKV(ctx, "serviceCommander.Delete: failed removing service", "err", err, "ID", serviceID)
		sendMessage(fmt.Sprintf("Fail to remove service with ID %d.", serviceID))
		return
	}

	sendMessage(fmt.Sprintf("Service with ID %d was deleted.", serviceID))
}
