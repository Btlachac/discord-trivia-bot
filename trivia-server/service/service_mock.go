// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	postgres "github.com/Btlachac/discord-trivia-bot/postgres"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddTrivia mocks base method.
func (m *MockService) AddTrivia(ctx *context.Context, newTrivia *postgres.Trivia) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTrivia", ctx, newTrivia)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTrivia indicates an expected call of AddTrivia.
func (mr *MockServiceMockRecorder) AddTrivia(ctx, newTrivia interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrivia", reflect.TypeOf((*MockService)(nil).AddTrivia), ctx, newTrivia)
}

// GetNewTrivia mocks base method.
func (m *MockService) GetNewTrivia() (*postgres.Trivia, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewTrivia")
	ret0, _ := ret[0].(*postgres.Trivia)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewTrivia indicates an expected call of GetNewTrivia.
func (mr *MockServiceMockRecorder) GetNewTrivia() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewTrivia", reflect.TypeOf((*MockService)(nil).GetNewTrivia))
}

// MarkTriviaUsed mocks base method.
func (m *MockService) MarkTriviaUsed(triviaId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkTriviaUsed", triviaId)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkTriviaUsed indicates an expected call of MarkTriviaUsed.
func (mr *MockServiceMockRecorder) MarkTriviaUsed(triviaId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkTriviaUsed", reflect.TypeOf((*MockService)(nil).MarkTriviaUsed), triviaId)
}

// RoundTypesList mocks base method.
func (m *MockService) RoundTypesList() ([]*postgres.RoundType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoundTypesList")
	ret0, _ := ret[0].([]*postgres.RoundType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RoundTypesList indicates an expected call of RoundTypesList.
func (mr *MockServiceMockRecorder) RoundTypesList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoundTypesList", reflect.TypeOf((*MockService)(nil).RoundTypesList))
}
