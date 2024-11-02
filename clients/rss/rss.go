package rss

import (
	"github.com/mmcdole/gofeed"
	"log"
	"net/url"
	"strings"
)

type ParsedNews struct {
	News     string
	NewsUrls string
}

func ValidateFeedURL(feedURL string) bool {

	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		log.Printf("invalid URL: %v", err)
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

func Parsing(feedURL string) ([]ParsedNews, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return []ParsedNews{}, err
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

		image := "Нет изображения"
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
				"- Заголовок: " + title +
				"\n  Image: " + image +
				"\n  Описание: " + description +
				"\n  Дата публикации: " + published +
				"\n  Автор: " + author +
				"\n  Ссылка: " + link,
				NewsUrls: link,
			})

		count++
	}

	return result, nil
}
