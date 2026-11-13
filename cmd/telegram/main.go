package main

import (
	"CNGDataManager/internal/config"
	"CNGDataManager/internal/telegram"
	"CNGDataManager/internal/utils"
	"log"

	"time"

	tele "gopkg.in/telebot.v4"
)

func main() {

	envConfig := config.NewEnvConfig()

	telegramBotConfig, err := config.NewTelegramBotConfig(envConfig)
	if err != nil {
		log.Fatal(err)
	}

	settings := tele.Settings{
		Token: telegramBotConfig.Token,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
		ParseMode: tele.ModeHTML,
		OnError:   telegram.OnErrorHandler,
	}

	apiClient := utils.NewAPIClient(telegramBotConfig.ApiBaseUrl)

	localStatesStorage := telegram.NewLocalStateHandlerStorage()

	c := telegram.NewHandlersController(telegramBotConfig, localStatesStorage, apiClient)

	bot, err := tele.NewBot(settings)

	if err != nil {
		log.Fatal(err)
	}

	bot.Use(telegram.AdminOnlyMiddleware(telegramBotConfig.AdminIds...))

	bot.Use(telegram.FSMMiddleware(localStatesStorage))

	localStatesStorage.Handle(telegram.StateStartDate, c.ReportHandlers.HandleStartDateInput)

	localStatesStorage.Handle(telegram.StateEndDate, c.ReportHandlers.HandleEndDateInput)

	bot.Handle("/start", c.StartHandler)

	bot.Handle(&c.Keyboards.MainMenuKeyboard.GenerateReportBtn, c.ReportHandlers.GenerateReportMenu)

	bot.Handle(&c.Keyboards.CreateReportKeyboard.BtnCreateWeekly, c.ReportHandlers.CreateWeeklyReport)

	bot.Handle(&c.Keyboards.CreateReportKeyboard.BtnCreateCurrentMonth, c.ReportHandlers.CreateCurrentMonthReport)

	bot.Handle(&c.Keyboards.CreateReportKeyboard.BtnCreateMonthly, c.ReportHandlers.CreateMonthlyReport)

	bot.Handle(&c.Keyboards.CreateReportKeyboard.BtnCreateFullTime, c.ReportHandlers.CreateFullTimeReport)

	bot.Handle(&c.Keyboards.CreateReportKeyboard.BtnCreateReportForOwnTime, c.ReportHandlers.CreateForOwnTimeReport)

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		return nil
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		return nil
	})

	bot.Start()
}
