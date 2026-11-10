package main

import (
	"CNGDataManager/internal/utils"
	"fmt"
	"log"
	"os"
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

	settings := tele.Settings{
		Token: botToken,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	bot, err := tele.NewBot(settings)

	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(c tele.Context) error {

		keyboard := tele.ReplyMarkup{

			ResizeKeyboard: true,
		}

		keyboard.Reply(
			keyboard.Row(
				keyboard.Text("Щотижневий звіт 📊"),
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

	btnCreateFullWeeklyReport := createRawWeeklyReportMenu.Data("Щотижневий звіт)", "weekly_report_full")

	btnCreateRawWeeklyReport := createRawWeeklyReportMenu.Data("Щотижневий звіт(цифри)", "weekly_report_raw")

	createRawWeeklyReportMenu.Inline(
		createRawWeeklyReportMenu.Row(btnCreateFullWeeklyReport),
		createRawWeeklyReportMenu.Row(btnCreateRawWeeklyReport),
	)

	bot.Handle("Щотижневий звіт 📊", func(c tele.Context) error {

		return c.Reply("Виберіть тип звіту для генерації... 👇",

			&tele.SendOptions{
				ParseMode:   tele.ModeHTML,
				ReplyMarkup: &createRawWeeklyReportMenu,
			})
	})

	bot.Handle(&btnCreateRawWeeklyReport, func(c tele.Context) error {

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

		e, th, lc, ms, _, _, err := utils.GetAllData(apiBase, params)

		if err != nil {
			return fmt.Errorf("error in getting all data: %v", err)
		}

		ed := utils.GetEvacuationBeneficiariesDisaggregation(e)

		thd := utils.GetTransitHubBeneficiariesDisaggregation(th)

		lcd := utils.GetTransportationsBeneficiariesDisaggregation(lc)

		msd := utils.GetMainShelterBeneficiariesDisaggregation(ms)

		if err := c.Edit(utils.BuildReportTextRaw(params.DateStart, params.DateEnd, ed, msd, thd, lcd)); err != nil {
			return fmt.Errorf("error in report generation: %v", err)
		}
		return nil
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

		allms, _, err := utils.GetMainShelterData(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all main shelter data: %v", err)
		}

		allTr, err := utils.GetTransitHubData(apiBase, utils.GetBeneficiariesDataParams{})

		if err != nil {
			return fmt.Errorf("error in getting all transit hub data: %v", err)
		}

		if err := c.Edit(utils.BuildFullReportText(params.DateStart, params.DateEnd, h, e, ms, allms, tr, allTr, lc, cUpdateTime),
			&tele.SendOptions{
				ParseMode: tele.ModeHTML,
			}); err != nil {
			return fmt.Errorf("error in report generation: %v", err)
		}
		return nil
	})

	bot.Start()
}
