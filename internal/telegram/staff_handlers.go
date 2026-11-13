package telegram

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

func CancelState(c tele.Context, session *StateSession) error {

	session.Clear()
	return c.EditOrReply("Дія скасована!")
}

func OnErrorHandler(err error, c tele.Context) {
	log.Printf("Telegram error: %v", err)

	if c != nil {
		_ = c.Send("Виникла помилка під час обробки запиту. Повідомте про це @KamikottoG")
	}
}
