package subscription

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/commands/subscription/service"
	"github.com/ozonmp/ssn-service-api/internal/bot/config"
	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/bot/path"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

const (
	serviceSudmn = "service"
)

type commander interface {
	HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(ctx context.Context, message *tgbotapi.Message, commandPath path.CommandPath)
}

type serviceService interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	List(ctx context.Context, offset uint64, limit uint64) (*subscription.ServiceListChunk, error)
	Create(ctx context.Context, service *subscription.Service) (uint64, error)
	Update(ctx context.Context, serviceID uint64, service *subscription.Service) error
	Remove(ctx context.Context, serviceID uint64) (bool, error)
}

type subscriptionCommander struct {
	bot              *tgbotapi.BotAPI
	serviceCommander commander
}

// NewSubscriptionCommander - creating subscription commander/
func NewSubscriptionCommander(bot *tgbotapi.BotAPI, srvService serviceService, cfg *config.Bot) *subscriptionCommander {
	return &subscriptionCommander{
		bot:              bot,
		serviceCommander: service.NewServiceCommander(bot, srvService, cfg.ListPerPage),
	}
}

func (c *subscriptionCommander) HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.Subdomain {
	case serviceSudmn:
		c.serviceCommander.HandleCallback(ctx, callback, callbackPath)
	default:
		logger.DebugKV(ctx, "subscriptionCommander.HandleCallback: unknown subdomain", "subdomain", callbackPath.Subdomain)
	}
}

func (c *subscriptionCommander) HandleCommand(ctx context.Context, msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.Subdomain {
	case serviceSudmn:
		c.serviceCommander.HandleCommand(ctx, msg, commandPath)
	default:
		logger.DebugKV(ctx, "subscriptionCommander.HandleCommand: unknown subdomain", "subdomain", commandPath.Subdomain)
	}
}
