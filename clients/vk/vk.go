package vk

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"telegramBot/lib/e"
	"time"
)

type Client struct {
	host    string
	version string
	token   string
	client  http.Client
}

func New(host, version, token string) *Client {
	return &Client{
		host:    host,
		version: version,
		token:   token,
		client:  http.Client{},
	}
}

func (c *Client) GetNews(groupID string) (news []string, err error) {
	defer func() { err = e.Wrap("vk-client/ can't get news", err) }()

	q := url.Values{}
	q.Add("count", "10")
	q.Add("domain", groupID)
	q.Add("access_token", c.token)
	q.Add("v", c.version)

	data, err := c.doRequest("wall.get", q)
	if err != nil {
		return nil, err
	}

	var res Response

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return vkParsing(res.Response.Items), nil
}

func (c *Client) ValidateNewsGroup(groupID string) (val bool, err error) {
	defer func() { err = e.Wrap("vk-client/ can't validate news group", err) }()

	q := url.Values{}
	q.Add("count", "1")
	q.Add("domain", groupID)
	q.Add("access_token", c.token)
	q.Add("v", c.version)

	data, err := c.doRequest("wall.get", q)
	if err != nil {
		return false, err
	}

	var res Response

	if err = json.Unmarshal(data, &res); err != nil {
		return false, err
	}

	return res.Response.Count > 0, nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.Wrap("vk-client/ can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join("method", method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func vkParsing(news []Post) []string {
	var result []string

	for _, item := range news {

		var media string

		for _, attachment := range item.Media {

			if attachment.Photo != nil {
				media += "\n  Photo: " + attachment.Photo.Sizes[0].URL
			}

			if attachment.Video != nil {
				media +=
					"\n  Video: " + attachment.Video.Image[0].URL +
						"\n  Описание видео: " + attachment.Video.Description
			}

			if attachment.Audio != nil {
				media +=
					"\n  Название аудио: " + attachment.Audio.Title +
						"\n  " + attachment.Audio.Artist +
						"\n  Audio: " + attachment.Audio.URL
			}

		}
		result = append(result,
			"-Заголовок: "+item.Text+
				media+
				"\n  Дата публикации: "+time.Unix(item.Date, 0).Format("2 January 2006 15:04"),
		)
	}

	return result
}
