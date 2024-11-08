package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"telegramBot/pkg/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(ctx context.Context, offset, limit int) (updates []Update, err error) {
	defer func() { err = e.Wrap("clients/telegram.Updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, "getUpdates", q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessageText(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("parse_mode", "HTML")

	_, err := c.doRequest(ctx, "sendMessage", q)
	if err != nil {
		return e.Wrap("clients/telegram.SendMessageText", err)
	}

	return nil
}

func (c *Client) SendMessageTextAndButton(ctx context.Context, chatID int, text string, button InlineKeyboardMarkup) (err error) {
	defer func() { err = e.Wrap("clients/telegram.SendMessageTextAndButton", err) }()

	jsonData, err := json.Marshal(button)
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("reply_markup", string(jsonData))
	q.Add("parse_mode", "HTML")

	_, err = c.doRequest(ctx, "sendMessage", q)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AnswerCallbackQuery(ctx context.Context, callbackID string, text string) error {
	q := url.Values{}
	q.Add("callback_query_id", callbackID)
	q.Add("text", text)

	_, err := c.doRequest(ctx, "answerCallbackQuery", q)
	if err != nil {
		return e.Wrap("clients/telegram.AnswerCallbackQuery", err)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.Wrap("clients/telegram.doRequest", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString(query.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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
