package l

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func LoggingStart(env string) (*slog.Logger, error) {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		file, err := createLogFile()
		if err != nil {
			return nil, err
		}
		log = slog.New(
			slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		file, err := createLogFile()
		if err != nil {
			return nil, err
		}
		log = slog.New(
			slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.Default()
	}

	return log, nil
}

func createLogFile() (*os.File, error) {
	logDir := "logs"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, 0774); err != nil {
			return nil, fmt.Errorf("can't create a logs dir: %w", err)
		}
	}

	nowDate := time.Now().Format(time.DateOnly)
	nowTime := strings.ReplaceAll(time.Now().Format(time.TimeOnly), ":", ".")

	file, err := os.Create(logDir + "/" + nowDate + "_" + nowTime + ".txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	return file, nil
}