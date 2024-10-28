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

type startObjects struct {
	jsonFilePath string
	tgToken      string
	vkToken      string
}

type jsonData struct {
	TgBotHost    string `json:"tgBotHost"`
	VkApiHost    string `json:"vkApiHost"`
	VkApiVersion string `json:"vkApiVersion"`
	ConnStr      string `json:"PSQLconnection"`
	BatchSize    int    `json:"batchSize"`
}

func main() {

	// logging
	if err := l.Start(); err != nil {
		log.Print(err)
	}

	startObjects := mustGetStartObjects()

	var launchData jsonData
	mustOpenJsonFiles(startObjects.jsonFilePath, &launchData)

	s, err := psql.New(launchData.ConnStr)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(launchData.TgBotHost, startObjects.tgToken),
		vkClient.New(launchData.VkApiHost, launchData.VkApiVersion, startObjects.vkToken),
		s,
	)

	go func() {
		log.Print("service started")

		consumer := eventconsumer.New(eventsProcessor, eventsProcessor, launchData.BatchSize)
		if err := consumer.Start(); err != nil {
			log.Fatal("service is stopped ", err)
		}
	}()

	// Program Completion
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Print("service stopping, sys signal: ", sign)

	eventconsumer.Stop()
}

func mustOpenJsonFiles(filePath string, launchData *jsonData) {
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

func mustGetStartObjects() *startObjects {

	var (
		jsonFilePath string
		tgToken      string
		vkToken      string
	)

	flag.StringVar(&jsonFilePath,
		"config-path",
		"",
		"path for access to config file",
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

	if jsonFilePath == "" {
		jsonFilePath = os.Getenv("CONFIG_PATH")
		if jsonFilePath == "" {
			log.Fatal("path is not specified")
		}
	}
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

	return &startObjects{
		jsonFilePath: jsonFilePath,
		tgToken:      tgToken,
		vkToken:      vkToken,
	}
}
