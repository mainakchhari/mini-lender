package errors

import "net/http"

type BaseError interface {
	Code() int
	Error() string
}

type InvalidPaymentAmountError struct{}

func (e InvalidPaymentAmountError) Error() string {
	return "payment amount cannot be less than due amount"
}

func (e InvalidPaymentAmountError) Code() int {
	return http.StatusNotAcceptable
}

type PaymentNotFoundError struct{}

func (e PaymentNotFoundError) Error() string {
	return "payment not found"
}

func (e PaymentNotFoundError) Code() int {
	return http.StatusNotFound
}
