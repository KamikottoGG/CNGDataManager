package telegram

import (
	"CNGDataManager/internal/config"
	"CNGDataManager/internal/utils"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v4"
)

type ReportHandlersController struct {
	Config       *config.TelegramBotConfig
	Keyboards    Keyboards
	StateStorage StateStorage
}

func NewReportHandlersController(config *config.TelegramBotConfig, keyboards Keyboards, storage StateStorage) *ReportHandlersController {
	return &ReportHandlersController{
		Config:       config,
		Keyboards:    keyboards,
		StateStorage: storage,
	}
}

func (h *ReportHandlersController) GenerateReportMenu(c tele.Context) error {

	return c.Reply("Виберіть тип звіту для генерації... 👇",

		&tele.SendOptions{
			ReplyMarkup: &h.Keyboards.CreateReportKeyboard.ReplyMarkup,
		})
}

func (h *ReportHandlersController) CreateWeeklyReport(c tele.Context) error {

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

	reportText, err := h.generateReportText(params)
	if err != nil {
		return err
	}
	if err := c.Edit(reportText); err != nil {
		return fmt.Errorf("error in report generation: %v", err)
	}
	return nil
}

func (h *ReportHandlersController) CreateCurrentMonthReport(c tele.Context) error {

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

	reportText, err := h.generateReportText(params)
	if err != nil {
		return err
	}
	if err := c.Edit(reportText); err != nil {
		return fmt.Errorf("error in report generation: %v", err)
	}
	return nil
}

func (h *ReportHandlersController) CreateMonthlyReport(c tele.Context) error {

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

	reportText, err := h.generateReportText(params)
	if err != nil {
		return err
	}
	if err := c.Edit(reportText); err != nil {
		return fmt.Errorf("error in report generation: %v", err)
	}
	return nil
}

func (h *ReportHandlersController) CreateFullTimeReport(c tele.Context) error {

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

	reportText, err := h.generateReportText(params)
	if err != nil {
		return err
	}
	if err := c.Edit(reportText); err != nil {
		return fmt.Errorf("error in report generation: %v", err)
	}
	return nil
}

func (h *ReportHandlersController) CreateForOwnTimeReport(c tele.Context) error {
	if err := c.Edit("Напишіть з якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> дд.мм.рррр \nНаприклад: <b><u>02.01.2006</u></b>",
		&tele.SendOptions{ReplyMarkup: h.Keyboards.NewCancelReplyMarkup(), ParseMode: tele.ModeHTML}); err != nil {
		return err
	}

	s := h.StateStorage.GetOrCreateSession(UserID(c.Sender().ID))

	s.SetState(StateStartDate)

	return nil
}

func (h *ReportHandlersController) HandleStartDateInput(c tele.Context, stateSession *StateSession) error {

	date, err := time.Parse("02.01.2006", c.Text())

	if err != nil {

		return c.Send("Увага! Дату не розпізнано.\nНапишіть з якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
			&tele.SendOptions{ParseMode: tele.ModeHTML})
	}

	if date.After(time.Now()) {
		return c.Send("Увага! Дата пізніше ніж сьогодні, перевірте будь ласка правильність написання дати\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
			&tele.SendOptions{ParseMode: tele.ModeHTML})
	}

	stateSession.Set("startDate", date)
	stateSession.SetState(StateEndDate)

	return c.Send("Напишіть ДО якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
		&tele.SendOptions{ParseMode: tele.ModeHTML})
}

func (h *ReportHandlersController) HandleEndDateInput(c tele.Context, stateSession *StateSession) error {

	startTimeStr := stateSession.Get("startDate")

	startTime, _ := startTimeStr.(time.Time)

	endTime, err := time.Parse("02.01.2006", c.Text())

	if err != nil {
		return c.Send("Увага! Дату не розпізнано.\nНапишіть ДО якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
			h.Keyboards.NewCancelReplyMarkup())
	}

	if startTime.After(endTime) {
		stateSession.Clear()

		return c.Send("Увага! Ви ввели некоректно дати, початкова дата("+startTime.Format("02.01.2006")+") більше ніж кінцева("+endTime.Format("02.01.2006")+")\n\nДія скасована.",
			&tele.SendOptions{ParseMode: tele.ModeHTML})
	}

	params := utils.GetBeneficiariesDataParams{
		DateStart: startTime,
		DateEnd:   endTime,
	}

	err = c.Send("Будуємо звіт для вас...")

	if err != nil {
		return err
	}

	stateSession.Clear()

	if err := c.Send(h.generateReportText(params)); err != nil {
		return fmt.Errorf("error in report generation: %v", err)
	}

	return nil
}

func (h *ReportHandlersController) generateReportText(params utils.GetBeneficiariesDataParams) (string, error) {
	e, tr, lc, ms, hot, cUpdateTime, err := utils.GetAllData(h.Config.ApiBaseUrl, params)

	if err != nil {
		return "", fmt.Errorf("error in getting all data: %v", err)
	}

	allms, _, err := utils.GetMainShelterRecords(h.Config.ApiBaseUrl, utils.GetBeneficiariesDataParams{})

	if err != nil {
		return "", fmt.Errorf("error in getting all main shelter data: %v", err)
	}

	allTr, err := utils.GetTransitHubRecords(h.Config.ApiBaseUrl, utils.GetBeneficiariesDataParams{})

	if err != nil {
		return "", fmt.Errorf("error in getting all transit hub data: %v", err)
	}

	vtr, err := utils.GetVehicleTripRecords(h.Config.ApiBaseUrl, params)

	if err != nil {
		return "", fmt.Errorf("error in getting vehicle trips data: %v", err)
	}

	fText := utils.BuildFullReportText(params.DateStart, params.DateEnd, hot, e, ms, allms, tr, allTr, lc, vtr, cUpdateTime)

	return fText, nil
}
