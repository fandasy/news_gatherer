package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"telegramBot/lib/e"
	"telegramBot/services"
	"telegramBot/storage"
)

const (
	HelpCmd        = "/help"
	StartCmd       = "/start"
	ListCmd        = "/list"
	AllNewsCmd     = "/allnews"
	NewsCertainCmd = "/news "
	RemoveCmd      = "/rm "

	VkGroupPath = "https://vk.com/"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	switch {
	case text == HelpCmd:
		return p.sendHelp(chatID)
	case text == StartCmd:
		return p.sendHello(chatID)
	case text == ListCmd:
		return p.sendList(ctx, chatID, username)
	case text == AllNewsCmd:
		return p.getAllNews(ctx, chatID, username)
	case strings.HasPrefix(text, NewsCertainCmd):
		return p.getСertainNews(ctx, chatID, username, text)
	case strings.HasPrefix(text, RemoveCmd):
		return p.removePage(ctx, chatID, text, username)
	case isAddCmd(text):
		// SaveNewsPages
		return p.defineAssembler(ctx, chatID, text, username)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) defineAssembler(ctx context.Context, chatID int, pageURL string, username string) error {
	switch {
	case strings.HasPrefix(pageURL, VkGroupPath):
		groupID := strings.TrimPrefix(pageURL, VkGroupPath)
		val, err := p.vk.ValidateNewsGroup(groupID)
		if err != nil {
			return err
		}
		if val {
			return p.savePage(ctx, chatID, pageURL, username, "VK")
		} else {
			return p.tg.SendMessage(chatID, msgNotValidateGroup)
		}
	case services.ValidateFeedURL(pageURL):
		return p.savePage(ctx, chatID, pageURL, username, "RSS")
	default:
		return p.tg.SendMessage(chatID, msgNotContainNewsFeed)
	}

	return nil
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string, assembler string) (err error) {
	defer func() { err = e.Wrap("can't do command: save page", err) }()

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
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) removePage(ctx context.Context, chatID int, rmPage string, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: remove page", err) }()

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
		return p.tg.SendMessage(chatID, msgNoSavedPagesRm+pageURL)
	}

	if err := p.storage.Remove(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgRemove); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendList(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: send list", err) }()

	pages, count, err := p.storage.PickPageList(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	msgList := generateListMsg(pages.URLS, count)

	if err := p.tg.SendMessage(chatID, msgList); err != nil {
		return err
	}

	return nil
}

func (p *Processor) getAllNews(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: get all news", err) }()

	newsList, err := p.storage.GetAllNews(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	for _, item := range newsList.News {

		newsArr, err := p.getNews(item)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, news := range newsArr {
			if err := p.tg.SendMessage(chatID, news); err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func (p *Processor) getСertainNews(ctx context.Context, chatID int, username string, cmdText string) (err error) {
	defer func() { err = e.Wrap("can't do command: get news", err) }()

	filter := strings.TrimPrefix(cmdText, NewsCertainCmd)

	switch {
	case strings.HasPrefix(filter, VkGroupPath):
		page := &storage.Page{
			URL:      filter,
			UserName: username,
		}

		isExists, err := p.storage.IsExists(ctx, page)
		if err != nil {
			return err
		}

		if !isExists {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		}

		newsArr, err := p.vk.GetNews(strings.TrimPrefix(filter, VkGroupPath))
		if err != nil {
			return err
		}

		for _, news := range newsArr {
			if err := p.tg.SendMessage(chatID, news); err != nil {
				log.Println(err)
			}
		}

	case filterIsAssembler(filter, "VK", "RSS"):
		page := &storage.Page{
			UserName:  username,
			Assembler: filter,
		}
		newsList, err := p.storage.PickNews(ctx, page)
		if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
			return err
		}

		if errors.Is(err, storage.ErrNoSavedPages) {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		}

		for _, item := range newsList.News {

			newsArr, err := p.getNews(item)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, news := range newsArr {
				if err := p.tg.SendMessage(chatID, news); err != nil {
					log.Println(err)
				}
			}
		}

	case services.ValidateFeedURL(filter):
		page := &storage.Page{
			URL:      filter,
			UserName: username,
		}

		isExists, err := p.storage.IsExists(ctx, page)
		if err != nil {
			return err
		}

		if !isExists {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		}

		newsArr, err := services.RSSParsing(filter)
		if err != nil {
			return err
		}

		for _, news := range newsArr {
			if err := p.tg.SendMessage(chatID, news); err != nil {
				log.Println(err)
			}
		}

	default:
		return p.tg.SendMessage(chatID, msgTypeOrPageNotExist)
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) getNews(item storage.News) ([]string, error) {
	switch item.Assembler {
	case "VK":
		newsArr, err := p.vk.GetNews(strings.TrimPrefix(item.URL, VkGroupPath))
		if err != nil {
			return nil, err
		}

		return newsArr, nil

	case "RSS":
		newsArr, err := services.RSSParsing(item.URL)
		if err != nil {
			return nil, err
		}

		return newsArr, nil
	}

	return nil, nil
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
