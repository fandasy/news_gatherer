package yagpt

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"telegramBot/lib/e"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	host   string
	token  string
	client *http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:   host,
		token:  newToken(token),
		client: &http.Client{},
	}
}

func newToken(token string) string {
	return "OAuth " + token
}

func (c *Client) GetRetelling(pageURL string) (retelling string, err error) {
	defer func() { err = e.Wrap("yaGPT-client/ can't get retelling", err) }()

	requestBody := RequestBody{
		ArticleURL: pageURL,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	data, err := c.doRequest("sharing-url", jsonData)
	if err != nil {
		return "", err
	}

	var res Response

	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	if res.Status == "error" {
		return "", ErrResponseStatusError
	}

	retelling, err = yaGptParsing(res.URL, pageURL)
	if err != nil {
		return "", err
	}

	return retelling, nil
}

func (c *Client) doRequest(method string, jsonData []byte) (data []byte, err error) {
	defer func() { err = e.Wrap("yaGPT-client/ can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join("api", method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func yaGptParsing(RetellingUrl, OriginalUrl string) (data string, err error) {
	defer func() { err = e.Wrap("yaGPT-client/ can't do parsing", err) }()

	res, err := http.Get(RetellingUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var result string

	summaryTitle := doc.Find(".summary-text .title").Text()
	result += summaryTitle + "\n"

	doc.Find(".summary-scroll .chapters .chapter").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".chapter-subheading").Text() // Заголовок
		result += "\nЗаголовок: " + title

		s.Find(".thesis").Each(func(j int, thesis *goquery.Selection) {
			thesisText := thesis.Find(".thesis-text").Text()
			result += thesisText + "\n"
		})
	})

	result += "\nСсылка на статью: " + OriginalUrl

	return result, nil
}
