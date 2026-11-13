package telegram

import tele "gopkg.in/telebot.v4"

func (h *HandlersController) StartHandler(c tele.Context) error {

	return c.Send("Вітаємо в боті менеджері CNG! 👋\n"+
		"\n"+
		"━━━━━━━━━━━━━━━━━━━━━\n"+
		"👇 Оберіть опцію, яку хочете виконати",
		&h.Keyboards.MainMenuKeyboard.ReplyKeyboard)
}
