package service

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/bot/path"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

const (
	helpCmd   = "help"
	listCmd   = "list"
	getCmd    = "get"
	deleteCmd = "delete"
	newCmd    = "new"
	editCmd   = "edit"
)

type serviceService interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	List(ctx context.Context, offset uint64, limit uint64) (*subscription.ServiceListChunk, error)
	Create(ctx context.Context, service *subscription.Service) (uint64, error)
	Update(ctx context.Context, serviceID uint64, service *subscription.Service) error
	Remove(ctx context.Context, serviceID uint64) (bool, error)
}

type serviceCommander struct {
	bot         *tgbotapi.BotAPI
	service     serviceService
	listPerPage uint64
}

// NewServiceCommander - creating service commander.
func NewServiceCommander(bot *tgbotapi.BotAPI, service serviceService, listPerPage uint64) *serviceCommander {
	return &serviceCommander{
		bot:         bot,
		service:     service,
		listPerPage: listPerPage,
	}
}

func (c *serviceCommander) HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case listCmd:
		c.CallbackList(ctx, callback, callbackPath)
	default:
		logger.DebugKV(ctx, "serviceCommander.HandleCallback: unknown callback name", "callback", callbackPath.CallbackName)
	}
}

func (c *serviceCommander) HandleCommand(ctx context.Context, msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.CommandName {
	case helpCmd:
		c.Help(ctx, msg)
	case getCmd:
		c.Get(ctx, msg)
	case listCmd:
		c.List(ctx, msg, 0, c.listPerPage)
	case deleteCmd:
		c.Delete(ctx, msg)
	case newCmd:
		c.New(ctx, msg)
	case editCmd:
		c.Edit(ctx, msg)
	default:
		c.Default(ctx, msg)
	}
}
