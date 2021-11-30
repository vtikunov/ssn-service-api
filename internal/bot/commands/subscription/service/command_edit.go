package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

func (c *serviceCommander) Edit(ctx context.Context, inputMessage *tgbotapi.Message) {
	args := inputMessage.CommandArguments()

	sendMessage := func(text string) {
		msg := tgbotapi.NewMessage(
			inputMessage.Chat.ID,
			text,
		)

		if _, err := c.bot.Send(msg); err != nil {
			logger.ErrorKV(ctx, "serviceCommander.Edit: error sending reply message to chat", "err", err)
		}
	}

	wrongArgs := func(args string) {
		logger.DebugKV(ctx, "serviceCommander.Edit: received wrong args", "args", args)
		sendMessage(fmt.Sprintf("Wrong arguments: \"%v\"! See: /help__subscription__service", args))
	}

	argsParts := strings.SplitN(args, "|", 3)
	if len(argsParts) != 3 {
		wrongArgs(args)
		return
	}

	serviceID, err := strconv.Atoi(argsParts[0])
	if err != nil {
		wrongArgs(args)
		return
	}

	name := argsParts[1]

	if len(name) == 0 {
		wrongArgs(args)
		return
	}

	desc := argsParts[2]

	if len(desc) == 0 {
		wrongArgs(args)
		return
	}

	serviceData := &subscription.Service{
		Name:        name,
		Description: desc,
	}

	if err := c.service.Update(ctx, uint64(serviceID), serviceData); err != nil {
		logger.ErrorKV(ctx, "serviceCommander.Update: failed updating service", "err", err, "ID", serviceID, "data", serviceData)
		sendMessage(fmt.Sprintf("Fail to update service with ID %d.", serviceID))
		return
	}

	sendMessage(fmt.Sprintf("Service with ID %d was updated.", serviceID))
}
