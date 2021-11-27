package service

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

func (c *serviceCommander) New(ctx context.Context, inputMessage *tgbotapi.Message) {
	args := inputMessage.CommandArguments()

	sendMessage := func(text string) {
		msg := tgbotapi.NewMessage(
			inputMessage.Chat.ID,
			text,
		)

		if _, err := c.bot.Send(msg); err != nil {
			logger.ErrorKV(ctx, "serviceCommander.New: error sending reply message to chat", "err", err)
		}
	}

	wrongArgs := func(args string) {
		logger.DebugKV(ctx, "serviceCommander.New: received wrong args", "args", args)
		sendMessage(fmt.Sprintf("Wrong arguments: \"%v\"! See: /help__subscription__service", args))
	}

	argsParts := strings.SplitN(args, "|", 2)
	if len(argsParts) != 2 {
		wrongArgs(args)
		return
	}

	name := argsParts[0]

	if len(name) == 0 {
		wrongArgs(args)
		return
	}

	desc := argsParts[1]

	if len(desc) == 0 {
		wrongArgs(args)
		return
	}

	serviceData := &subscription.Service{
		Name:        name,
		Description: desc,
	}

	serviceID, err := c.service.Create(ctx, serviceData)

	if err != nil {
		logger.ErrorKV(ctx, "serviceCommander.New: failed creating service", "err", err)
		sendMessage("Fail to create service.")
		return
	}

	sendMessage(fmt.Sprintf("Service with ID %d was created.", serviceID))
}
