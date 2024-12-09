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
	"telegramBot/internal/events/telegram"
	"telegramBot/internal/lib/logger/l"
	"telegramBot/internal/storage/psql"
)

func main() {

	flagsObj, err := s.GetFlagsObjects()
	if err != nil {
		panic(err)
	}

	cfg, err := j.LoadConfig(flagsObj.JsonFilePath)
	if err != nil {
		panic(err)
	}

	log, err := l.LoggingStart(cfg.Env)
	if err != nil {
		panic(err)
	}

	storage, err := psql.New(cfg.ConnStr, log)
	if err != nil {
		panic(err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		panic(err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(cfg.Clients.TgBotHost, flagsObj.TgToken),
		vkClient.New(cfg.Clients.VkApiHost, cfg.Clients.VkApiVersion, flagsObj.VkToken),
		yaGptClient.New(cfg.Clients.YaGptHost, flagsObj.YaGptToken),
		storage,
		cfg.ReqLimit,
		log,
	)

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, cfg.BatchSize, log)

	log.Info("service started")

	go consumer.Start(cfg.UpdateTimeout)

	// Program Completion
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("service stopping", slog.Any("sys signal", sign))

	consumer.Stop()
}
