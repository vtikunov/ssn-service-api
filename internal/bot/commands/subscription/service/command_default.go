package service

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

func (c *serviceCommander) Default(ctx context.Context, inputMessage *tgbotapi.Message) {
	logger.DebugKV(ctx, "serviceCommander.Default: received message", "from", inputMessage.From.UserName, "text", inputMessage.Text)

	var sb strings.Builder
	sb.WriteString("You wrote: ")
	sb.WriteString(inputMessage.Text)

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, sb.String())

	if _, err := c.bot.Send(msg); err != nil {
		logger.ErrorKV(ctx, "serviceCommander.Default: error sending reply message to chat", "err", err)
	}
}
