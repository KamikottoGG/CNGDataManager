package utils

import tele "gopkg.in/telebot.v4"

func AdminOnlyMiddleware(adminIDs ...int64) tele.MiddlewareFunc {
	admins := make(map[int64]struct{}, len(adminIDs))
	for _, id := range adminIDs {
		admins[id] = struct{}{}
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if _, ok := admins[c.Sender().ID]; !ok {

				return c.Send("Нажаль, у вас відсутній доступ до цієї функції.")
			}
			return next(c)
		}
	}
}
