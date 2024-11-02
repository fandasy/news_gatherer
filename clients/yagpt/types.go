package yagpt

import "errors"

var (
	ErrResponseStatusError = errors.New("yaGPT response status error")
)

type Response struct {
	Status string `json:"status"`
	URL    string `json:"sharing_url"`
}

type RequestBody struct {
	ArticleURL string `json:"article_url"`
}
