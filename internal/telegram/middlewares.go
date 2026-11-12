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

type FSMState int

const (
	STATE_START_DATE FSMState = iota
	STATE_END_DATE
)

func FSMMiddleware(stateStorage map[int64]FSMState, stateDataStorage map[int64]map[string]string) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			userID := c.Sender().ID

			callback := c.Callback()

			if callback != nil {
				cdata := strings.TrimSpace(c.Callback().Data)

				if cdata == "cancel" {
					return CancelState(c, stateDataStorage, stateStorage)
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

			state, hasState := stateStorage[userID]
			if !hasState {
				return next(c)
			}

			if stateDataStorage[userID] == nil {
				stateDataStorage[userID] = make(map[string]string)
			}

			switch state {
			case STATE_START_DATE:

				if c.Message() != nil {
					return HandleStartDateInput(c, stateDataStorage[userID], stateStorage)
				}
			case STATE_END_DATE:
				if c.Message() != nil {
					return HandleEndDateInput(c, stateDataStorage[userID], stateStorage)
				}
			default:
				delete(stateStorage, userID)
				delete(stateDataStorage, userID)
			}

			return next(c)
		}
	}
}
