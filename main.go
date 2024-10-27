package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgClient "telegramBot/clients/telegram"
	vkClient "telegramBot/clients/vk"
	"telegramBot/consumer/event-consumer"
	"telegramBot/events/telegram"
	"telegramBot/lib/l"
	"telegramBot/storage/psql"
)

type BotToken struct {
	TgToken string
	VkToken string
}

type JSONData struct {
	TgBotHost    string `json:"tgBotHost"`
	VkApiHost    string `json:"vkApiHost"`
	VkApiVersion string `json:"vkApiVersion"`
	ConnStr      string `json:"PSQLconnection"`
	BatchSize    int    `json:"batchSize"`
}

const jsonFilePath = "data.json"

// batchSize - updatesBatchLimit, between 1 - 100, defaults to 100

func main() {

	// logs
	if err := l.Start(); err != nil {
		log.Print(err)
	}

	var launchData JSONData

	openJSONfiles(jsonFilePath, &launchData)

	s, err := psql.New(launchData.ConnStr)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	tokens := mustToken()

	eventsProcessor := telegram.New(
		tgClient.New(launchData.TgBotHost, tokens.TgToken),
		vkClient.New(launchData.VkApiHost, launchData.VkApiVersion, tokens.VkToken),
		s,
	)

	go func() {
		log.Print("service started")

		consumer := eventconsumer.New(eventsProcessor, eventsProcessor, launchData.BatchSize)
		if err := consumer.Start(); err != nil {
			log.Fatal("service is stopped ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Print("service stopping, sys signal: ", sign)

	eventconsumer.Stop()
}

func openJSONfiles(filePath string, launchData *JSONData) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Failed to open JSON file: ", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Failed to read JSON file: ", err)
	}

	if err := json.Unmarshal(byteValue, launchData); err != nil {
		log.Fatal("Failed to parse JSON file: ", err)
	}
}

func mustToken() *BotToken {

	var (
		tgToken string
		vkToken string
	)

	flag.StringVar(&tgToken,
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.StringVar(&vkToken,
		"vk-bot-token",
		"",
		"token for access to vk bot",
	)
	flag.Parse()

	if tgToken == "" {
		tgToken = os.Getenv("TG_TOKEN")
		if tgToken == "" {
			log.Fatal("tgToken is not specified")
		}
	}

	if vkToken == "" {
		vkToken = os.Getenv("VK_TOKEN")
		if vkToken == "" {
			log.Fatal("vkToken is not specified")
		}
	}

	return &BotToken{
		TgToken: tgToken,
		VkToken: vkToken,
	}
}
