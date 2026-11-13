package telegram

import "CNGDataManager/internal/config"

type HandlersController struct {
	Config         *config.TelegramBotConfig
	Keyboards      Keyboards
	ReportHandlers *ReportHandlersController
	StateStorage   StateStorage
}

func NewHandlersController(config *config.TelegramBotConfig, storage StateStorage) *HandlersController {
	c := HandlersController{
		Config:       config,
		Keyboards:    NewKeyboards(),
		StateStorage: storage,
	}

	c.ReportHandlers = NewReportHandlersController(c.Config, c.Keyboards, storage)

	return &c
}
