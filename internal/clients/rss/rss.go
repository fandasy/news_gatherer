package rss

import (
	"context"
	"github.com/mmcdole/gofeed"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"telegramBot/internal/lib/logger/sl"
	"telegramBot/pkg/e"
)

type ParsedNews struct {
	News     string
	NewsUrls string
}

func ValidateFeedURL(ctx context.Context, feedURL string) bool {
	const op = "rss.ValidateFeedURL"

	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		slog.Error(op, sl.Err(err))
		return false
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		slog.Error(op, sl.Err(err))
		return false
	}

	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		slog.Error(op, sl.Err(err))
		return false
	}

	if resp.StatusCode != 200 {
		slog.Info(op, slog.Int("StatusCode", resp.StatusCode))
		return false
	}

	if strings.HasSuffix(parsedURL.Path, ".rss") || strings.HasSuffix(parsedURL.Path, ".xml") {
		return true
	}

	if strings.Contains(strings.ToLower(parsedURL.Path), "rss") || strings.Contains(strings.ToLower(parsedURL.Path), "feed") {
		return true
	}

	return false
}

func Parsing(ctx context.Context, feedURL string) ([]ParsedNews, error) {
	const op = "rss.Parsing"

	defer func() {
		if err := recover(); err != nil {
			slog.Error("[RECOVER] panic err: ", err)
		}
	}()

	fp := gofeed.NewParser()

	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return []ParsedNews{}, e.Wrap(op, err)
	}

	var (
		result []ParsedNews
		count  int
	)

	for _, item := range feed.Items {
		if count == 10 {
			break
		}

		title := "Не указано"
		if item.Title != "" {
			title = item.Title
		}

		image := ""
		if item.Image != nil && item.Image.URL != "" {
			image = item.Image.URL
		}

		description := "Нет описания"
		if item.Description != "" {
			description = item.Description
		}

		published := "Дата не указана"
		if item.Published != "" {
			published = item.Published
		}

		author := "Автор не указан"
		if item.Author != nil && item.Author.Name != "" {
			author = item.Author.Name
		}

		link := item.Link

		result = append(result,
			ParsedNews{News: "" +
				title + "\n" +
				`<a href="` + image + `"> Image </a>` +
				"\nОписание: " + description +
				"\nДата публикации: " + published +
				"\nАвтор: " + author + "\n" +
				`<a href="` + link + `"> Ссылка на статью</a>`,
				NewsUrls: link,
			})

		count++
	}

	return result, nil
}
