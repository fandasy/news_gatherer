package j

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"telegramBot/pkg/e"
	"time"
)

type ConfigData struct {
	Env           string        `yaml:"env"`
	Clients       Clients       `yaml:"clients"`
	ConnStr       string        `yaml:"PSQLConnection"`
	BatchSize     int           `yaml:"batchSize"`
	UpdateTimeout time.Duration `yaml:"updateTimeout"`
	ReqLimit      ReqLimit      `yaml:"reqLimit"`
}

type Clients struct {
	TgBotHost    string `yaml:"tgBotHost"`
	VkApiHost    string `yaml:"vkApiHost"`
	VkApiVersion string `yaml:"vkApiVersion"`
	YaGptHost    string `yaml:"yaGptHost"`
}

type ReqLimit struct {
	MaxNumberReq uint32        `yaml:"maxNumberReq"`
	TimeSlice    time.Duration `yaml:"timeSlice"`
	BanTime      time.Duration `yaml:"banTime"`
}

func LoadConfig(filePath string) (*ConfigData, error) {
	var cfg ConfigData

	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("failed to open JSON file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, e.Wrap("failed to read JSON file", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, e.Wrap("failed to parse JSON file", err)
	}

	return &cfg, err
}
