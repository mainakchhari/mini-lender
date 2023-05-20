// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/domain/factory/payment.go

// Package mock_factory is a generated GoMock package.
package mock_factory

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/mainakchhari/mini-lender/internal/app/domain"
)

// MockIPayment is a mock of IPayment interface.
type MockIPayment struct {
	ctrl     *gomock.Controller
	recorder *MockIPaymentMockRecorder
}

// MockIPaymentMockRecorder is the mock recorder for MockIPayment.
type MockIPaymentMockRecorder struct {
	mock *MockIPayment
}

// NewMockIPayment creates a new mock instance.
func NewMockIPayment(ctrl *gomock.Controller) *MockIPayment {
	mock := &MockIPayment{ctrl: ctrl}
	mock.recorder = &MockIPaymentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPayment) EXPECT() *MockIPaymentMockRecorder {
	return m.recorder
}

// GeneratePaymentsFromLoan mocks base method.
func (m *MockIPayment) GeneratePaymentsFromLoan(loanEnt domain.Loan, NumInstalments int) []domain.Payment {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratePaymentsFromLoan", loanEnt, NumInstalments)
	ret0, _ := ret[0].([]domain.Payment)
	return ret0
}

// GeneratePaymentsFromLoan indicates an expected call of GeneratePaymentsFromLoan.
func (mr *MockIPaymentMockRecorder) GeneratePaymentsFromLoan(loanEnt, NumInstalments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratePaymentsFromLoan", reflect.TypeOf((*MockIPayment)(nil).GeneratePaymentsFromLoan), loanEnt, NumInstalments)
}
