package telegram

import (
	"CNGDataManager/internal/utils"
	"fmt"
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v4"
)

func HandleStartDateInput(c tele.Context, stateData map[string]string, stateStorage map[int64]FSMState) error {

	date, err := time.Parse("02.01.2006", c.Text())

	if err != nil {

		return c.Send("Увага! Дату не роспізнано.\nНапишіть з якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
			&tele.SendOptions{ParseMode: tele.ModeHTML})
	}

	if date.After(time.Now()) {
		return c.Send("Увага! Дата пізніше ніж сьогодні, перевірте будь ласка правильність написання дати\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
			&tele.SendOptions{ParseMode: tele.ModeHTML})
	}

	stateData["startDate"] = date.Format("02.01.2006")
	stateStorage[c.Sender().ID] = STATE_END_DATE

	return c.Send("Напишіть ДО якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
		&tele.SendOptions{ParseMode: tele.ModeHTML})
}

func HandleEndDateInput(c tele.Context, stateData map[string]string, stateStorage map[int64]FSMState) error {

	startTimeStr := stateData["startDate"]

	startTime, _ := time.Parse("02.01.2006", startTimeStr)

	endTime, err := time.Parse("02.01.2006", c.Text())

	if err != nil {
		return c.Send("Увага! Дату не роспізнано.\nНапишіть ДО якої дати потрібно знегерувати звіт.\n\nУВАГА! Формат <b>обов'язкого</b> <b><u>дд.мм.рррр</u></b> \nНаприклад: <b><u>02.01.2006</u></b>",
			&tele.SendOptions{ParseMode: tele.ModeHTML})
	}

	if startTime.After(endTime) {
		delete(stateData, "startDate")
		delete(stateStorage, c.Sender().ID)
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

	delete(stateData, "startDate")
	delete(stateStorage, c.Sender().ID)

	apiBase := os.Getenv("API_BASE_URL")

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

	if err := c.Send(utils.BuildFullReportText(params.DateStart, params.DateEnd, h, e, ms, allms, tr, allTr, lc, vtr, cUpdateTime),
		&tele.SendOptions{
			ParseMode: tele.ModeHTML,
		}); err != nil {
		return fmt.Errorf("error in report generation: %v", err)
	}

	return nil
}

func CancelState(c tele.Context, statesData map[int64]map[string]string, statesStorage map[int64]FSMState) error {

	delete(statesStorage, c.Sender().ID)
	delete(statesData, c.Sender().ID)
	return c.Respond(&tele.CallbackResponse{Text: "Дія скасована!"})
}

func OnErrorHandler(err error, c tele.Context) {
	log.Printf("Telegram error: %v", err)

	if c != nil {
		_ = c.Send("Виникла помилка під час обробки запиту. Повідомте про це @KamikottoG")
	}
}
