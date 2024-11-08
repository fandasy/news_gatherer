package s

import (
	"flag"
	"fmt"
	"os"
)

type StartObjects struct {
	JsonFilePath string
	TgToken      string
	VkToken      string
	YaGptToken   string
}

func GetFlagsObjects() (*StartObjects, error) {

	var (
		jsonFilePath string
		tgToken      string
		vkToken      string
		yaGptToken   string
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
	flag.StringVar(&yaGptToken,
		"ya-gpt-token",
		"",
		"token for access to yaGpt token",
	)
	flag.Parse()

	if jsonFilePath == "" {
		jsonFilePath = os.Getenv("CONFIG_PATH")
		if jsonFilePath == "" {
			return nil, fmt.Errorf("config path is not specified")
		}
	}
	if tgToken == "" {
		tgToken = os.Getenv("TG_TOKEN")
		if tgToken == "" {
			return nil, fmt.Errorf("tgToken is not specified")
		}
	}
	if vkToken == "" {
		vkToken = os.Getenv("VK_TOKEN")
		if vkToken == "" {
			return nil, fmt.Errorf("vkToken is not specified")
		}
	}
	if yaGptToken == "" {
		yaGptToken = os.Getenv("YA_GPT_TOKEN")
		if yaGptToken == "" {
			return nil, fmt.Errorf("yaGptToken is not specified")
		}
	}

	return &StartObjects{
		JsonFilePath: jsonFilePath,
		TgToken:      tgToken,
		VkToken:      vkToken,
		YaGptToken:   yaGptToken,
	}, nil
}
