package services

import (
	"github.com/mmcdole/gofeed"
	"log"
	"net/url"
	"strings"
)

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

func RSSParsing(feedURL string) ([]string, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return []string{}, err
	}

	var (
		result []string
		count  int
	)

	for _, item := range feed.Items {
		if count == 10 {
			break
		}

		result = append(result, `
- Заголовок: `+item.Title+`
  Ссылка: `+item.Link+`
  Описание: `+item.Description+`
  Дата публикации: `+item.Published+`
  Автор: `+item.Author.Name+`
  Image: `+item.Image.URL)

		count++
	}

	return result, nil
}
