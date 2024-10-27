package l

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func Start() error {
	logDir := "logs"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, 0774); err != nil {
			return fmt.Errorf("can't create a logs file: %w", err)
		}
	}

	nowDate := time.Now().Format(time.DateOnly)
	nowTime := strings.ReplaceAll(time.Now().Format(time.TimeOnly), ":", ".")

	file, err := os.Create(logDir + "/" + nowDate + "_" + nowTime + ".txt")
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	return nil
}
