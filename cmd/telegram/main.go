package main

import (
	"CNGDataManager/internal/telegram"
	"CNGDataManager/internal/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

func main() {

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	apiBase := os.Getenv("API_BASE_URL")

	if apiBase == "" {
		log.Fatal("API_BASE_URL environment variable not set")
	}

	adminIds := os.Getenv("ADMIN_IDS")

	if adminIds == "" {
		log.Fatal("ADMIN_IDS environment variable not set")
	}

	adminIdsSplit := strings.Split(adminIds, ",")

	adminIdsInt := make([]int64, len(adminIdsSplit))

	for _, s := range adminIdsSplit {
		id, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		adminIdsInt = append(adminIdsInt, id)
	}

	settings := tele.Settings{
		Token: botToken,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
		OnError: telegram.OnErrorHandler,
	}

	bot, err := tele.NewBot(settings)

	if err != nil {
		log.Fatal(err)
	}

	bot.Use(telegram.AdminOnlyMiddleware(adminIdsInt...))

	states := make(map[int64]telegram.FSMState)
	statesData := make(map[int64]map[string]string)

	bot.Use(telegram.FSMMiddleware(states, statesData))

	bot.Handle("/start", func(c tele.Context) error {

		keyboard := tele.ReplyMarkup{

			ResizeKeyboard: true,
		}

		keyboard.Reply(
			keyboard.Row(
				keyboard.Text("Сгенерувати звіт 📊"),
			),
		)

		return c.Send("Вітаємо в боті менеджері CNG! 👋\n"+
			"\n"+
			"━━━━━━━━━━━━━━━━━━━━━\n"+
			"👇 Оберіть опцію, яку хочете виконати",
			&keyboard,
		)
	})

	createRawWeeklyReportMenu := tele.ReplyMarkup{}

	btnCreateFullWeeklyReport := createRawWeeklyReportMenu.Data("Щотижневий звіт", "weekly_report_full")
	btnCreateFullMontlyReport := createRawWeeklyReportMenu.Data("Щомісячний звіт", "monthly_report_full")
	btnCreateFullCurrentMonthReport := createRawWeeklyReportMenu.Data("Цей місяць", "current_month_report_full")
	btnCreateFullReport := createRawWeeklyReportMenu.Data("Увесь час", "all_report_full")
	btnCreateOwnReport := createRawWeeklyReportMenu.Data("Свій проміжок", "own_report_full")

	//btnCreateRawWeeklyReport := createRawWeeklyReportMenu.Data("Сгенерувати звіт(малий)", "weekly_report_raw")

	createRawWeeklyReportMenu.Inline(
		createRawWeeklyReportMenu.Row(btnCreateFullWeeklyReport),
		createRawWeeklyReportMenu.Row(btnCreateFullMontlyReport),
		createRawWeeklyReportMenu.Row(btnCreateFullCurrentMonthReport),
		createRawWeeklyReportMenu.Row(btnCreateFullReport),
		createRawWeeklyReportMenu.Row(btnCreateOwnReport),
	)

	bot.Handle("Сгенерувати звіт 📊", func(c tele.Context) error {

		return c.Reply("Виберіть тип звіту для генерації... 👇",

			&tele.SendOptions{
				ParseMode:   tele.ModeHTML,
				ReplyMarkup: &createRawWeeklyReportMenu,
			})
	})

	bot.Handle(&btnCreateFullWeeklyReport, func(c tele.Context) error {

		if err := c.Respond(&tele.CallbackResponse{}); err != nil {
			return err
		}

		if err := c.Edit("Генерую звіт..."); err != nil {
			return err
		}

		now := time.Now()

		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		currentMonday := now.AddDate(0, 0, -weekday+1)

		lastMonday := currentMonday.AddDate(0, 0, -7)

		lastSunday := lastMonday.AddDate(0, 0, 6)

		lastMonday = time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, now.Location())
		lastSunday = time.Date(lastSunday.Year(), lastSunday.Month(), lastSunday.Day(), 0, 0, 0, 0, now.Location())

		params := utils.GetBeneficiariesDataParams{
			DateStart: lastMonday,
			DateEnd:   lastSunday,
		}

		e, tr, lc, ms, h, cUpdateTime, err := utils.GetAllData(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting all data: %v", err)
		}

		allms, _, err := utils.GetMainShelterRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all main shelter data: %v", err)
		}

		allTr, err := utils.GetTransitHubRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all transit hub data: %v", err)
		}

		vtr, err := utils.GetVehicleTripRecords(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting vehicle trips data: %v", err)
		}

		if err := c.Edit(utils.BuildFullReportText(params.DateStart, params.DateEnd, h, e, ms, allms, tr, allTr, lc, vtr, cUpdateTime),
			&tele.SendOptions{
				ParseMode: tele.ModeHTML,
			}); err != nil {
			return fmt.Errorf("error in report generation: %v", err)
		}
		return nil
	})

	bot.Handle(&btnCreateFullCurrentMonthReport, func(c tele.Context) error {

		if err := c.Respond(&tele.CallbackResponse{}); err != nil {
			return err
		}

		if err := c.Edit("Генерую звіт..."); err != nil {
			return err
		}

		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		params := utils.GetBeneficiariesDataParams{
			DateStart: firstOfMonth,
			DateEnd:   lastOfMonth,
		}

		e, tr, lc, ms, h, cUpdateTime, err := utils.GetAllData(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting all data: %v", err)
		}

		allms, _, err := utils.GetMainShelterRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all main shelter data: %v", err)
		}

		allTr, err := utils.GetTransitHubRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all transit hub data: %v", err)
		}

		vtr, err := utils.GetVehicleTripRecords(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting vehicle trips data: %v", err)
		}

		if err := c.Edit(utils.BuildFullReportText(params.DateStart, params.DateEnd, h, e, ms, allms, tr, allTr, lc, vtr, cUpdateTime),
			&tele.SendOptions{
				ParseMode: tele.ModeHTML,
			}); err != nil {
			return fmt.Errorf("error in report generation: %v", err)
		}
		return nil
	})

	bot.Handle(&btnCreateFullMontlyReport, func(c tele.Context) error {

		if err := c.Respond(&tele.CallbackResponse{}); err != nil {
			return err
		}

		if err := c.Edit("Генерую звіт..."); err != nil {
			return err
		}

		lastMonth := time.Now().AddDate(0, -1, 0)
		currentYear, currentMonth, _ := lastMonth.Date()
		currentLocation := lastMonth.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		params := utils.GetBeneficiariesDataParams{
			DateStart: firstOfMonth,
			DateEnd:   lastOfMonth,
		}

		e, tr, lc, ms, h, cUpdateTime, err := utils.GetAllData(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting all data: %v", err)
		}

		allms, _, err := utils.GetMainShelterRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all main shelter data: %v", err)
		}

		allTr, err := utils.GetTransitHubRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all transit hub data: %v", err)
		}

		vtr, err := utils.GetVehicleTripRecords(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting vehicle trips data: %v", err)
		}

		if err := c.Edit(utils.BuildFullReportText(params.DateStart, params.DateEnd, h, e, ms, allms, tr, allTr, lc, vtr, cUpdateTime),
			&tele.SendOptions{
				ParseMode: tele.ModeHTML,
			}); err != nil {
			return fmt.Errorf("error in report generation: %v", err)
		}
		return nil
	})

	bot.Handle(&btnCreateFullReport, func(c tele.Context) error {

		if err := c.Respond(&tele.CallbackResponse{}); err != nil {
			return err
		}

		if err := c.Edit("Генерую звіт..."); err != nil {
			return err
		}

		params := utils.GetBeneficiariesDataParams{
			DateStart: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().Location()),
			DateEnd:   time.Now(),
		}

		e, tr, lc, ms, h, cUpdateTime, err := utils.GetAllData(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting all data: %v", err)
		}

		allms, _, err := utils.GetMainShelterRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all main shelter data: %v", err)
		}

		allTr, err := utils.GetTransitHubRecords(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all transit hub data: %v", err)
		}

		vtr, err := utils.GetVehicleTripRecords(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting vehicle trips data: %v", err)
		}

		if err := c.Edit(utils.BuildFullReportText(params.DateStart, params.DateEnd, h, e, ms, allms, tr, allTr, lc, vtr, cUpdateTime),
			&tele.SendOptions{
				ParseMode: tele.ModeHTML,
			}); err != nil {
			return fmt.Errorf("error in report generation: %v", err)
		}
		return nil
	})

	cancelMenu := tele.ReplyMarkup{}

	cancelBtn := cancelMenu.Data("Скасувати дію", "cancel")

	cancelMenu.Inline(
		cancelMenu.Row(
			cancelBtn),
	)

	bot.Handle(&btnCreateOwnReport, func(c tele.Context) error {

		if err := c.Respond(); err != nil {
			return err
		}

		if err = c.Edit("Напишіть з якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> дд.мм.рррр \nНаприклад: <b><u>02.01.2006</u></b>",
			&tele.SendOptions{ReplyMarkup: &cancelMenu, ParseMode: tele.ModeHTML}); err != nil {
			return err
		}

		states[c.Sender().ID] = telegram.STATE_START_DATE

		return nil
	})

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		return nil
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		return nil
	})

	bot.Start()
}
