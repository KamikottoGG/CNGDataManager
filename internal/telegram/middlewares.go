package telegram

import (
	"strings"

	tele "gopkg.in/telebot.v4"
)

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

func FSMMiddleware(stateStorage StateStorage) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {

			userID := c.Sender().ID

			userSession := stateStorage.GetOrCreateSession(
				UserID(userID),
			)

			callback := c.Callback()

			if callback != nil {
				cdata := strings.TrimSpace(c.Callback().Data)

				if cdata == "cancel" {
					return CancelState(c, userSession)
				}
			}

			//isCallbackCancel := callback != nil && strings.TrimSpace(c.Callback().Unique) == "cancel"
			//
			//isTextCancel := c.Message() != nil && (c.Text() == "cancel" || c.Text() == "/cancel")
			//
			//if isCallbackCancel || isTextCancel {
			//
			//	return CancelState(c, stateDataStorage, stateStorage)
			//}

			if userSession.GetState() == StateNone {
				return next(c)
			}

			fun, err := stateStorage.GetHandler(userSession.GetState())

			if err != nil {
				return err
			}

			return fun(c, userSession)
		}
	}
}
