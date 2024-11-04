package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"telegramBot/lib/j"
	"telegramBot/lib/s"

	tgClient "telegramBot/clients/telegram"
	vkClient "telegramBot/clients/vk"
	yaGptClient "telegramBot/clients/yagpt"
	"telegramBot/consumer/event-consumer"
	"telegramBot/events/telegram"
	"telegramBot/lib/l"
	"telegramBot/storage/psql"
)

func main() {

	startObjects, err := s.MustGetStartObjects()
	if err != nil {
		panic(err)
	}

	launchData, err := j.MustOpenJsonFiles(startObjects.JsonFilePath)
	if err != nil {
		panic(err)
	}

	log, err := l.LoggingStart(launchData.Env)
	if err != nil {
		panic(err)
	}

	storage, err := psql.New(launchData.ConnStr, log)
	if err != nil {
		panic(err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		panic(err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(launchData.TgBotHost, startObjects.TgToken),
		vkClient.New(launchData.VkApiHost, launchData.VkApiVersion, startObjects.VkToken),
		yaGptClient.New(launchData.YaGptHost, startObjects.YaGptToken),
		storage,
		log,
	)

	go func() {
		log.Info("service started")

		consumer := eventconsumer.New(eventsProcessor, eventsProcessor, launchData.BatchSize, log)

		consumer.Start()

		log.Info("service is stopped")
	}()

	// Program Completion
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("service stopping", slog.Any("sys signal", sign))

	eventconsumer.Stop()
}
