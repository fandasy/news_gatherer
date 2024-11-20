package telegram

import (
	"context"
	"errors"
	"log/slog"
	telegram2 "telegramBot/internal/clients/telegram"
	"telegramBot/internal/clients/vk"
	"telegramBot/internal/clients/yagpt"
	"telegramBot/internal/config/j"
	"telegramBot/internal/controller/req-controller"
	"telegramBot/internal/events"
	"telegramBot/internal/storage"
	"telegramBot/pkg/e"
	"telegramBot/pkg/shortener"
)

type Processor struct {
	tg         *telegram2.Client
	vk         *vk.Client
	yaGpt      *yagpt.Client
	offset     int
	storage    storage.Storage
	log        *slog.Logger
	urlsMap    *shortener.UrlsMap
	reqCounter *req_controller.ReqCounter
}

type Meta struct {
	ChatID          int
	Username        string
	CallbackQueryID string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(

	tgClient *telegram2.Client,
	vkClient *vk.Client,
	yaGptClient *yagpt.Client,
	storage storage.Storage,
	reqLimitOptions j.ReqLimit,
	log *slog.Logger,

) *Processor {

	return &Processor{
		tg:         tgClient,
		vk:         vkClient,
		yaGpt:      yaGptClient,
		storage:    storage,
		log:        log,
		urlsMap:    shortener.NewUrlsMap(),
		reqCounter: req_controller.New(reqLimitOptions),
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {

	updates, err := p.tg.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	const op = "telegram/telegram.Process: "

	if event.Type == events.Message || event.Type == events.Callback {
		return p.processMessage(ctx, event)
	}

	return e.Wrap(op, ErrUnknownEventType)
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	const op = "telegram/telegram.processMessage: "

	meta, err := meta(event)
	if err != nil {
		return e.Wrap(op, err)
	}

	if err := p.doCmd(ctx, event.Text, meta); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	const op = "telegram/telegram.meta: "

	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap(op, ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram2.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	if updType == events.Callback {
		res.Meta = Meta{
			ChatID:          upd.CallbackQuery.Message.Chat.ID,
			Username:        upd.CallbackQuery.Message.From.Username,
			CallbackQueryID: upd.CallbackQuery.ID,
		}
	}

	return res
}

func fetchText(upd telegram2.Update) string {

	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}

	if upd.Message != nil {
		return upd.Message.Text
	}

	return ""
}

func fetchType(upd telegram2.Update) int {
	if upd.CallbackQuery != nil {
		return events.Callback
	}
	if upd.Message != nil {
		return events.Message
	}

	return events.Unknown
}
