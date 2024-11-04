package telegram

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strings"
	"telegramBot/clients/rss"
	"telegramBot/clients/telegram"
	"telegramBot/lib/e"
	"telegramBot/lib/sl"
	"telegramBot/storage"
)

const (
	HelpCmd         = "/help"
	StartCmd        = "/start"
	ListCmd         = "/list"
	AllNewsCmd      = "/allnews"
	ConcreteNewsCmd = "/news "
	RemoveCmd       = "/rm "
	RetellingCmd    = "/retelling "

	VkGroupPath = "https://vk.com/"

	maxMessageSize = 4096
)

func (p *Processor) doCmd(ctx context.Context, text string, meta Meta) error {
	text = strings.TrimSpace(text)

	chatID := meta.ChatID
	username := meta.Username
	callbackID := meta.CallbackQueryID

	log := p.log.With(
		slog.String("command", text),
		slog.String("username", username),
	)

	log.Info("got new command")

	switch {
	case text == HelpCmd:
		return p.sendHelp(ctx, chatID)

	case text == StartCmd:
		return p.sendHello(ctx, chatID)

	case text == ListCmd:
		return p.sendList(ctx, chatID, username)

	case text == AllNewsCmd:
		return p.getAllNews(ctx, chatID, username)

	case strings.HasPrefix(text, ConcreteNewsCmd):
		return p.getConcreteNews(ctx, chatID, username, text)

	case strings.HasPrefix(text, RemoveCmd):
		return p.removePage(ctx, chatID, text, username)

	case strings.HasPrefix(text, RetellingCmd):
		return p.retelling(ctx, chatID, callbackID, text)

	case isAddCmd(text):
		// SaveNewsPages
		return p.defineAssembler(ctx, chatID, text, username)

	default:
		return p.tg.SendMessageText(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) defineAssembler(ctx context.Context, chatID int, pageURL string, username string) error {
	switch {
	case strings.HasPrefix(pageURL, VkGroupPath):
		groupID := strings.TrimPrefix(pageURL, VkGroupPath)
		val, err := p.vk.ValidateNewsGroup(ctx, groupID)
		if err != nil {
			return err
		}
		if val {
			return p.savePage(ctx, chatID, pageURL, username, "VK")
		} else {
			return p.tg.SendMessageText(ctx, chatID, msgNotValidateGroup)
		}

	case rss.ValidateFeedURL(pageURL):
		return p.savePage(ctx, chatID, pageURL, username, "RSS")

	default:
		return p.tg.SendMessageText(ctx, chatID, msgNotContainNewsFeed)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string, assembler string) (err error) {
	defer func() { err = e.Wrap("commands/telegram.savePage:", err) }()

	page := &storage.Page{
		URL:       pageURL,
		UserName:  username,
		Assembler: assembler,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessageText(ctx, chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessageText(ctx, chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) removePage(ctx context.Context, chatID int, rmPage string, username string) (err error) {
	defer func() { err = e.Wrap("commands/telegram.removePage", err) }()

	pageURL := strings.TrimPrefix(rmPage, RemoveCmd)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}

	if !isExists {
		return p.tg.SendMessageText(ctx, chatID, msgNoSavedPagesRm+pageURL)
	}

	if err := p.storage.Remove(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessageText(ctx, chatID, msgRemove); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendList(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("commands/telegram.sendList", err) }()

	pages, count, err := p.storage.PickPageList(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessageText(ctx, chatID, msgNoSavedPages)
	}

	msgList := generateListMsg(pages.URLS, count)

	if err := p.tg.SendMessageText(ctx, chatID, msgList); err != nil {
		return err
	}

	return nil
}

func (p *Processor) retelling(ctx context.Context, chatID int, callbackID string, text string) (err error) {
	defer func() { err = e.Wrap("commands/telegram.retelling", err) }()

	shortKey := strings.TrimPrefix(text, RetellingCmd)

	originalUrl, ok := p.urlsMap.Get(shortKey)
	if !ok {
		// Так как команда /retelling может приниматься как текст с сообщения, приложение будет пропускать обработку ссылок не из мапы
		return p.tg.SendMessageText(ctx, chatID, msgImpossibleRetelling)
	}

	// Ответ на запрос обратного вызова, что бы кнопка прекратила мерцать
	if err := p.tg.AnswerCallbackQuery(ctx, callbackID, msgRetellingStarted); err != nil {
		return err
	}

	retelling, err := p.yaGpt.GetRetelling(ctx, originalUrl)
	if err != nil {
		return err
	}

	if len(retelling) > maxMessageSize {
		retellingArr := splitMessage(retelling)

		for _, retelling := range retellingArr {
			if err := p.tg.SendMessageText(ctx, chatID, retelling); err != nil {
				return err
			}
		}

	} else {

		if err := p.tg.SendMessageText(ctx, chatID, retelling); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) getAllNews(ctx context.Context, chatID int, username string) error {

	const op = "commands/telegram.getAllNews"

	log := p.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	newsFeedList, err := p.storage.GetAllNews(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return e.Wrap(op, err)
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessageText(ctx, chatID, msgNoSavedPages)
	}

	for _, newsFeedInfo := range newsFeedList.News {
		if err := p.getNewsAndSendMessage(ctx, chatID, newsFeedInfo); err != nil {
			log.Error("", sl.Err(err))
		}
	}

	return nil
}

func (p *Processor) getConcreteNews(ctx context.Context, chatID int, username string, cmdText string) (err error) {
	const op = "commands/telegram.getConcreteNews"

	log := p.log.With(
		slog.String("op", "commands/telegram.getConcreteNews"),
	)

	filter := strings.TrimPrefix(cmdText, ConcreteNewsCmd)

	if filterIsAssembler(filter, "VK", "RSS") {

		page := &storage.Page{
			UserName:  username,
			Assembler: filter,
		}
		newsFeedList, err := p.storage.PickNews(ctx, page)
		if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
			return e.Wrap(op, err)
		}

		if errors.Is(err, storage.ErrNoSavedPages) {
			return p.tg.SendMessageText(ctx, chatID, msgNoSavedPages)
		}

		for _, newsFeedInfo := range newsFeedList.News {
			if err := p.getNewsAndSendMessage(ctx, chatID, newsFeedInfo); err != nil {
				log.Error("", sl.Err(err))
			}
		}

		return nil
	}

	page := &storage.Page{
		URL:      filter,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return e.Wrap(op, err)
	}

	if !isExists {
		return p.tg.SendMessageText(ctx, chatID, msgNoSavedPages)
	}

	switch {
	case strings.HasPrefix(filter, VkGroupPath):

		newsFeedInfo := storage.News{
			URL:       filter,
			Assembler: "VK",
		}

		if err := p.getNewsAndSendMessage(ctx, chatID, newsFeedInfo); err != nil {
			return e.Wrap(op, err)
		}

	case rss.ValidateFeedURL(filter):

		newsFeedInfo := storage.News{
			URL:       filter,
			Assembler: "RSS",
		}

		if err := p.getNewsAndSendMessage(ctx, chatID, newsFeedInfo); err != nil {
			return e.Wrap(op, err)
		}

	default:
		return p.tg.SendMessageText(ctx, chatID, msgTypeOrPageNotExist)
	}

	return nil
}

func (p *Processor) getNewsAndSendMessage(ctx context.Context, chatID int, newsFeedInfo storage.News) (err error) {
	defer func() { err = e.Wrap("commands/telegram.getNewsAndSendMessage", err) }()

	switch newsFeedInfo.Assembler {
	case "VK":
		parsedNewsArr, err := p.vk.GetNews(ctx, strings.TrimPrefix(newsFeedInfo.URL, VkGroupPath))
		if err != nil {
			return err
		}

		for _, news := range parsedNewsArr {
			if err := p.tg.SendMessageText(ctx, chatID, news); err != nil {
				return err
			}
		}

	case "RSS":
		parsedNewsArr, err := rss.Parsing(ctx, newsFeedInfo.URL)
		if err != nil {
			return err
		}

		for _, parsedNews := range parsedNewsArr {

			// Для CallbackData есть ограничение на размер 64 байта, из-за этого нужно генерировать для каждого URL свою сокращённую версию
			shortKey := p.urlsMap.GenerateShortKey(parsedNews.NewsUrls)

			button := telegram.InlineKeyboardMarkup{
				InlineKeyboard: [][]telegram.InlineKeyboardButton{
					{
						{
							Text:         "Пересказать",
							CallbackData: "/retelling " + shortKey,
						},
					},
				},
			}

			if err := p.tg.SendMessageTextAndButton(ctx, chatID, parsedNews.News, button); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessageText(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessageText(ctx, chatID, msgHello)
}

func isAddCmd(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

func filterIsAssembler(filter string, assembler ...string) bool {
	for _, as := range assembler {
		if filter == as {
			return true
		}
	}

	return false
}

func splitMessage(s string) []string {
	var result []string

	for len(s) > maxMessageSize {
		result = append(result, s[:maxMessageSize])
		s = s[maxMessageSize:]
	}

	if len(s) > 0 {
		result = append(result, s)
	}

	return result
}
