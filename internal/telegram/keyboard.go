package telegram

import (
	tele "gopkg.in/telebot.v4"
)

type Keyboards struct {
	CreateReportKeyboard CreateReportKeyboard
	CancelInlineBtn      tele.Btn
	MainMenuKeyboard     MainMenuKeyboard
}

func NewKeyboards() Keyboards {
	kb := Keyboards{}

	kb.CreateReportKeyboard = newCreateReportKeyboard()

	kb.CancelInlineBtn = tele.Btn{
		Unique: "cancel",
		Text:   "Скасувати дію",
	}

	kb.MainMenuKeyboard = newMainMenuKeyboard()

	return kb
}

func (k *Keyboards) NewCancelReplyMarkup() *tele.ReplyMarkup {
	kb := tele.ReplyMarkup{}

	kb.Inline(
		kb.Row(k.CancelInlineBtn))

	return &kb
}

//btnCreateFullWeeklyReport := createRawWeeklyReportMenu.Data("Щотижневий звіт", "weekly_report_full")
//btnCreateFullMontlyReport := createRawWeeklyReportMenu.Data("Щомісячний звіт", "monthly_report_full")
//btnCreateFullCurrentMonthReport := createRawWeeklyReportMenu.Data("Цей місяць", "current_month_report_full")
//btnCreateFullReport := createRawWeeklyReportMenu.Data("Увесь час", "all_report_full")
//btnCreateOwnReport := createRawWeeklyReportMenu.Data("Свій проміжок", "own_report_full")

type MainMenuKeyboard struct {
	ReplyKeyboard     tele.ReplyMarkup
	GenerateReportBtn tele.Btn
}

func newMainMenuKeyboard() MainMenuKeyboard {
	kb := tele.ReplyMarkup{
		ResizeKeyboard: true,
	}
	btn := kb.Text("Сгенерувати звіт 📊")

	kb.Reply(
		kb.Row(
			btn,
		),
	)

	return MainMenuKeyboard{
		ReplyKeyboard:     kb,
		GenerateReportBtn: btn,
	}
}

type CreateReportKeyboard struct {
	ReplyMarkup               tele.ReplyMarkup
	BtnCreateWeekly           tele.Btn
	BtnCreateMonthly          tele.Btn
	BtnCreateCurrentMonth     tele.Btn
	BtnCreateFullTime         tele.Btn
	BtnCreateReportForOwnTime tele.Btn
}

func newCreateReportKeyboard() CreateReportKeyboard {

	kb := CreateReportKeyboard{
		ReplyMarkup: tele.ReplyMarkup{},
	}

	kb.BtnCreateWeekly = kb.ReplyMarkup.Data("Щотижневий звіт", "weekly_report_full")
	kb.BtnCreateMonthly = kb.ReplyMarkup.Data("Щомісячний звіт", "monthly_report_full")
	kb.BtnCreateCurrentMonth = kb.ReplyMarkup.Data("Цей місяць", "current_month_report_full")
	kb.BtnCreateFullTime = kb.ReplyMarkup.Data("Увесь час", "all_report_full")
	kb.BtnCreateReportForOwnTime = kb.ReplyMarkup.Data("Свій проміжок", "own_report_full")

	kb.ReplyMarkup.Inline(
		kb.ReplyMarkup.Row(
			kb.BtnCreateWeekly,
		),
		kb.ReplyMarkup.Row(
			kb.BtnCreateMonthly,
		),

		kb.ReplyMarkup.Row(
			kb.BtnCreateCurrentMonth,
		),
		kb.ReplyMarkup.Row(
			kb.BtnCreateFullTime,
		),
		kb.ReplyMarkup.Row(
			kb.BtnCreateReportForOwnTime,
		),
	)

	return kb
}
