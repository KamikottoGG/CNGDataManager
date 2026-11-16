package telegram

import (
	"CNGDataManager/internal/utils"
	"strconv"

	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

type DataGuardController struct {
	ApiClient    BeneficiaryAPI
	Keyboards    Keyboards
	StateStorage StateStorage
}

func newDataGuardController(api BeneficiaryAPI, kb Keyboards, ss StateStorage) *DataGuardController {
	return &DataGuardController{
		ApiClient:    api,
		Keyboards:    kb,
		StateStorage: ss,
	}
}

func (d *DataGuardController) HandleDataGuardMenuButton(c tele.Context) error {

	return c.Send("Вітаємо в меню DataGuard, тут можно слідкувати за якістю данних та їх своєчасним оновленням. \n\nДля продовження виберите опцію з меню нижче 👇",
		&d.Keyboards.DataGuardMenuKeyboard.ReplyMarkup)
}

func (d *DataGuardController) HandleDataGuardCheckLastRecords(c tele.Context) error {

	if err := c.Respond(); err != nil {
		return err
	}

	repData, err := d.ApiClient.GetAllData(utils.GetBeneficiariesDataParams{})

	if err != nil {
		return err
	}

	sBuilder := strings.Builder{}
	now := time.Now()

	sBuilder.WriteString("<b>Звіт по останнім зафіксованим показникам</b>\n")
	sBuilder.WriteString("<b>Дані актуальні на: <u>")
	sBuilder.WriteString(repData.CacheUpdateTime.Format("02.01.2006 15:04:05"))
	sBuilder.WriteString("</u></b>\n\n")

	wSection := func(direction string, monitoringIndicator string, lastDate time.Time) {
		sBuilder.WriteString("Напрямок: ")
		sBuilder.WriteString("<b><u>")
		sBuilder.WriteString(direction)
		sBuilder.WriteString("</u></b>")

		sBuilder.WriteByte('\n')

		sBuilder.WriteString("Моніторинг по індикатору: ")
		sBuilder.WriteString("<b>")
		sBuilder.WriteString(monitoringIndicator)
		sBuilder.WriteString("</b>")
		sBuilder.WriteByte('\n')

		sBuilder.WriteString("Остання фіксація: ")

		if lastDate.Day() == now.Day() && lastDate.Month() == now.Month() && lastDate.Year() == now.Year() {
			sBuilder.WriteString("<b><u>")
			sBuilder.WriteString("Сьогодні")
			sBuilder.WriteString("</u></b>")
			sBuilder.WriteByte('\n')
			sBuilder.WriteByte('\n')
			return
		}
		sBuilder.WriteString(lastDate.Format("02.01.2006"))
		sBuilder.WriteString(" <b><u> (")
		dn := now.Day() - lastDate.Day()
		sBuilder.WriteString(strconv.Itoa(dn))
		if dn == 1 {
			sBuilder.WriteString(" день назад)")
		}
		if dn >= 2 && dn <= 5 {
			sBuilder.WriteString(" дні назад)")
		}
		if dn > 5 {
			sBuilder.WriteString(" днів назад)")
		}

		sBuilder.WriteString("</u></b>")

		sBuilder.WriteByte('\n')
		sBuilder.WriteByte('\n')
	}

	wSection("Евакуація", "Дата евакуації", utils.GetLastIndicator(repData.Evacuation))
	wSection("Локальні перевезення", "Дата перевезення", utils.GetLastIndicator(repData.LocalTransportations))
	wSection("Гаряча лінія", "Дата дзвінка", utils.GetLastIndicator(repData.Hotline))
	wSection("Транзитний хаб", "Дата заїзду", utils.GetLastIndicator(repData.TransitHub))
	wSection("Шелтер", "Дата заїзду", utils.GetLastIndicator(repData.Shelter))
	wSection("Звіт з авто", "Дата запису", utils.GetLastIndicator(repData.VehicleTrip))

	return c.EditOrReply(sBuilder.String())

}
