package telegram

import (
	tele "gopkg.in/telebot.v4"
)

type Keyboards struct {
	CreateReportKeyboard  CreateReportKeyboard
	DataGuardMenuKeyboard DataGuardKeyboard
	CancelInlineBtn       tele.Btn
	MainMenuKeyboard      MainMenuKeyboard
}

func NewKeyboards() Keyboards {
	kb := Keyboards{}

	kb.CreateReportKeyboard = newCreateReportKeyboard()

	kb.CancelInlineBtn = tele.Btn{
		Unique: "cancel",
		Text:   "Скасувати дію",
	}

	kb.MainMenuKeyboard = newMainMenuKeyboard()

	kb.DataGuardMenuKeyboard = newDataGuardKeyboard()

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
	DataGuardMenuBtn  tele.Btn
}

func newMainMenuKeyboard() MainMenuKeyboard {
	kb := tele.ReplyMarkup{
		ResizeKeyboard: true,
	}
	rBtn := kb.Text("Сгенерувати звіт 📊")
	gBtn := kb.Text("Data Guard 🛡️")
	kb.Reply(
		kb.Row(
			rBtn,
		),
		kb.Row(
			gBtn,
		),
	)

	return MainMenuKeyboard{
		ReplyKeyboard:     kb,
		GenerateReportBtn: rBtn,
		DataGuardMenuBtn:  gBtn,
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

type DataGuardKeyboard struct {
	ReplyMarkup  tele.ReplyMarkup
	BtnCheckLast tele.Btn
}

func newDataGuardKeyboard() DataGuardKeyboard {
	d := DataGuardKeyboard{}

	d.ReplyMarkup = tele.ReplyMarkup{}
	d.BtnCheckLast = d.ReplyMarkup.Data("Останні записи", "check_last_records")
	d.ReplyMarkup.Inline(
		d.ReplyMarkup.Row(
			d.BtnCheckLast,
		),
	)

	return d
}
