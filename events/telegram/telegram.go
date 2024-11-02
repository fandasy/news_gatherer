package telegram

import (
	"context"
	"errors"

	"telegramBot/clients/telegram"
	"telegramBot/clients/vk"
	"telegramBot/clients/yagpt"
	"telegramBot/events"
	"telegramBot/lib/e"
	"telegramBot/lib/shortener"
	"telegramBot/storage"
)

type Processor struct {
	tg      *telegram.Client
	vk      *vk.Client
	yaGpt   *yagpt.Client
	offset  int
	storage storage.Storage
	urlsMap *shortener.UrlsMap
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

func New(tgClient *telegram.Client, vkClient *vk.Client, yaGptClient *yagpt.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      tgClient,
		vk:      vkClient,
		yaGpt:   yaGptClient,
		storage: storage,
		urlsMap: shortener.NewUrlsMap(),
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {

	updates, err := p.tg.Updates(p.offset, limit)
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

	if event.Type == events.Message || event.Type == events.Callback {
		return p.processMessage(ctx, event)
	}

	return e.Wrap("can't process message", ErrUnknownEventType)
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
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

func fetchText(upd telegram.Update) string {

	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}

	if upd.Message != nil {
		return upd.Message.Text
	}

	return ""
}

func fetchType(upd telegram.Update) events.Type {
	if upd.CallbackQuery != nil {
		return events.Callback
	}
	if upd.Message != nil {
		return events.Message
	}

	return events.Unknown
}
