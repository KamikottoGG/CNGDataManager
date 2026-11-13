package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	telegramBotTokenEnv = "TELEGRAM_BOT_TOKEN"
	adminIdsEnv         = "ADMIN_IDS"
	apiBaseUrl          = "API_BASE_URL"
)

type TelegramBotConfig struct {
	Token      string
	AdminIds   []int64
	ApiBaseUrl string
}

func NewTelegramBotConfig(config Config) (*TelegramBotConfig, error) {

	tc := TelegramBotConfig{}

	telegramBotToken := config.Get(telegramBotTokenEnv)

	if telegramBotToken == "" {
		return nil, errors.New("telegram bot token required " + telegramBotTokenEnv)
	}
	tc.Token = telegramBotToken
	adminIdsStr := config.Get(adminIdsEnv)
	if adminIdsStr == "" {
		return nil, errors.New("admin ids required " + adminIdsEnv)
	}

	adminIds := strings.Split(adminIdsStr, ",")

	tc.AdminIds = []int64{}

	for _, id := range adminIds {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing admin ids %s: %s", id, err)
		}
		tc.AdminIds = append(tc.AdminIds, idInt)
	}

	apiBase := config.Get(apiBaseUrl)
	if apiBase == "" {
		return nil, errors.New("api base url required " + apiBaseUrl)
	}
	tc.ApiBaseUrl = apiBase

	return &tc, nil

}
