package service

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

func (c *serviceCommander) Help(ctx context.Context, inputMessage *tgbotapi.Message) {
	var sb strings.Builder

	sb.WriteString("/help__subscription__service — print list of commands\n")
	sb.WriteString("/get__subscription__service id — get a service\n")
	sb.WriteString("/list__subscription__service — get a list of services\n")
	sb.WriteString("/delete__subscription__service id — delete an existing service\n")
	sb.WriteString("/new__subscription__service Name|Description — create a new service\n")
	sb.WriteString("/edit__subscription__service ID|Name|Description — edit a service")

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, sb.String())

	if _, err := c.bot.Send(msg); err != nil {
		logger.ErrorKV(ctx, "serviceCommander.Help: error sending reply message to chat", "err", err)
	}
}
