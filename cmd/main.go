package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	tgClient "telegramBot/internal/clients/telegram"
	vkClient "telegramBot/internal/clients/vk"
	yaGptClient "telegramBot/internal/clients/yagpt"
	"telegramBot/internal/config/j"
	"telegramBot/internal/config/s"
	"telegramBot/internal/consumer/event-consumer"
	req_controller "telegramBot/internal/controller/req-controller"
	"telegramBot/internal/events/telegram"
	"telegramBot/internal/lib/logger/l"
	"telegramBot/internal/storage/psql"
)

func main() {

	startObjects, err := s.GetFlagsObjects()
	if err != nil {
		panic(err)
	}

	launchData, err := j.OpenJsonFiles(startObjects.JsonFilePath)
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

	reqLimit := req_controller.NewLimitOptions(launchData.MaxNumberReq, launchData.TimeSlice, launchData.BanTime)

	eventsProcessor := telegram.New(
		tgClient.New(launchData.TgBotHost, startObjects.TgToken),
		vkClient.New(launchData.VkApiHost, launchData.VkApiVersion, startObjects.VkToken),
		yaGptClient.New(launchData.YaGptHost, startObjects.YaGptToken),
		storage,
		reqLimit,
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
