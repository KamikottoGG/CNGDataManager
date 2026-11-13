package telegram

import (
	"errors"
	"sync"

	"gopkg.in/telebot.v4"
)

type FSMState string

type StateHandler func(telebot.Context, *StateSession) error

const (
	StateNone      FSMState = "none"
	StateStartDate          = "start_date"
	StateEndDate            = "end_date"
)

type UserID int64

type StateHandlerStorage interface {
	Handle(FSMState, StateHandler)
	GetHandler(FSMState) (StateHandler, error)
}

type LocalStateHandlersStorage struct {
	data map[FSMState]StateHandler
}

func (l *LocalStateHandlersStorage) Handle(state FSMState, handler StateHandler) {
	l.data[state] = handler
}

func (l *LocalStateHandlersStorage) GetHandler(state FSMState) (StateHandler, error) {
	handler, ok := l.data[state]
	if !ok {

		return nil, errors.New("state not found")
	}
	return handler, nil
}

type StateSession struct {
	currentState FSMState

	data map[string]any

	mutex sync.RWMutex
}

func (s *StateSession) SetState(ns FSMState) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.currentState = ns

}

func (s *StateSession) GetState() FSMState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.currentState
}

func (s *StateSession) GetData() map[string]any {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data
}

func (s *StateSession) SetData(data map[string]any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = data
}

func (s *StateSession) Set(k string, v any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[k] = v
}
func (s *StateSession) Get(k string) any {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[k]
}
func (s *StateSession) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.currentState = StateNone
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}

	for _, k := range keys {
		delete(s.data, k)
	}

}

type StateStorage interface {
	GetOrCreateSession(userID UserID) *StateSession
	StateHandlerStorage
}

type LocalStateStorage struct {
	usersState map[UserID]*StateSession
	LocalStateHandlersStorage

	mutex sync.RWMutex
}

func (l *LocalStateStorage) GetOrCreateSession(userID UserID) *StateSession {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	ss, ok := l.usersState[userID]

	if !ok {
		l.usersState[userID] = &StateSession{
			currentState: StateNone,
			data:         make(map[string]any),
		}
		return l.usersState[userID]
	}

	return ss
}

func NewLocalStateHandlerStorage() *LocalStateStorage {

	return &LocalStateStorage{
		usersState: make(map[UserID]*StateSession),
		data:       make(map[FSMState]StateHandler),
	}

}
