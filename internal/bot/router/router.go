package router

import (
	"context"
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/config"
	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/bot/path"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	ssn_command "github.com/ozonmp/ssn-service-api/internal/bot/commands/subscription"
)

const (
	subscriptionDmn = "subscription"
)

type commander interface {
	HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(ctx context.Context, callback *tgbotapi.Message, commandPath path.CommandPath)
}

type serviceService interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	List(ctx context.Context, offset uint64, limit uint64) (*subscription.ServiceListChunk, error)
	Create(ctx context.Context, service *subscription.Service) (uint64, error)
	Update(ctx context.Context, serviceID uint64, service *subscription.Service) error
	Remove(ctx context.Context, serviceID uint64) (bool, error)
}

type router struct {
	bot                   *tgbotapi.BotAPI
	subscriptionCommander commander
}

// NewRouter - creating router.
func NewRouter(bot *tgbotapi.BotAPI, srvService serviceService, cfg *config.Bot) *router {
	return &router{
		bot:                   bot,
		subscriptionCommander: ssn_command.NewSubscriptionCommander(bot, srvService, cfg),
	}
}

func (c *router) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			logger.ErrorKV(ctx, "recovered from panic", "panicValue", panicValue, "stack", string(debug.Stack()))
		}
	}()

	switch {
	case update.CallbackQuery != nil:
		c.handleCallback(ctx, update.CallbackQuery)
	case update.Message != nil:
		c.handleMessage(ctx, update.Message)
	}
}

func (c *router) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	callbackPath, err := path.ParseCallback(callback.Data)
	if err != nil {
		logger.ErrorKV(ctx, "router.handleCallback: failed parsing callback data", "err", err, "data", callback.Data)

		return
	}

	switch callbackPath.Domain {
	case subscriptionDmn:
		c.subscriptionCommander.HandleCallback(ctx, callback, callbackPath)
	default:
		logger.DebugKV(ctx, "router.handleCallback: unknown domain", "domain", callbackPath.Domain)
	}
}

func (c *router) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		c.showCommandFormat(ctx, msg)

		return
	}

	commandPath, err := path.ParseCommand(msg.Command())
	if err != nil {
		logger.ErrorKV(ctx, "router.handleMessage: failed parsing message data", "err", err, "data", msg.Command())
		return
	}

	switch commandPath.Domain {
	case subscriptionDmn:
		c.subscriptionCommander.HandleCommand(ctx, msg, commandPath)
	default:
		logger.DebugKV(ctx, "router.handleMessage: unknown domain", "domain", commandPath.Domain)
	}
}

func (c *router) showCommandFormat(ctx context.Context, inputMessage *tgbotapi.Message) {
	outputMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Command format: /{command}__{domain}__{subdomain}")

	_, err := c.bot.Send(outputMsg)
	if err != nil {
		logger.ErrorKV(ctx, "router.showCommandFormat: error sending reply message to chat", "err", err)
	}
}
