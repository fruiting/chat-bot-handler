// Code generated by MockGen. DO NOT EDIT.
// Source: chat_bot_processor.go

// Package internal is a generated GoMock package.
package internal

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockChatBotHandler is a mock of ChatBotHandler interface.
type MockChatBotHandler struct {
	ctrl     *gomock.Controller
	recorder *MockChatBotHandlerMockRecorder
}

// MockChatBotHandlerMockRecorder is the mock recorder for MockChatBotHandler.
type MockChatBotHandlerMockRecorder struct {
	mock *MockChatBotHandler
}

// NewMockChatBotHandler creates a new mock instance.
func NewMockChatBotHandler(ctrl *gomock.Controller) *MockChatBotHandler {
	mock := &MockChatBotHandler{ctrl: ctrl}
	mock.recorder = &MockChatBotHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatBotHandler) EXPECT() *MockChatBotHandlerMockRecorder {
	return m.recorder
}

// FindChatIdAndText mocks base method.
func (m *MockChatBotHandler) FindChatIdAndText(bodyRequest []byte) (ChatId, Text, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindChatIdAndText", bodyRequest)
	ret0, _ := ret[0].(ChatId)
	ret1, _ := ret[1].(Text)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindChatIdAndText indicates an expected call of FindChatIdAndText.
func (mr *MockChatBotHandlerMockRecorder) FindChatIdAndText(bodyRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindChatIdAndText", reflect.TypeOf((*MockChatBotHandler)(nil).FindChatIdAndText), bodyRequest)
}

// SendMessage mocks base method.
func (m *MockChatBotHandler) SendMessage(chatId ChatId, text string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", chatId, text)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockChatBotHandlerMockRecorder) SendMessage(chatId, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockChatBotHandler)(nil).SendMessage), chatId, text)
}
