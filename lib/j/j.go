package j

import (
	"encoding/json"
	"io"
	"os"
	"telegramBot/lib/e"
)

type JsonData struct {
	Env          string `json:"env"`
	TgBotHost    string `json:"tgBotHost"`
	VkApiHost    string `json:"vkApiHost"`
	VkApiVersion string `json:"vkApiVersion"`
	YaGptHost    string `json:"yaGptHost"`
	ConnStr      string `json:"PSQLconnection"`
	BatchSize    int    `json:"batchSize"`
}

func MustOpenJsonFiles(filePath string) (*JsonData, error) {
	var launchData JsonData

	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("failed to open JSON file", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, e.Wrap("failed to read JSON file", err)
	}

	if err := json.Unmarshal(byteValue, &launchData); err != nil {
		return nil, e.Wrap("failed to parse JSON file", err)
	}

	return &launchData, err
}
