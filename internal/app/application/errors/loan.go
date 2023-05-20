package errors

import "net/http"

type LoanNotFoundError struct{}

func (e LoanNotFoundError) Error() string {
	return "loan not found"
}

func (e LoanNotFoundError) Code() int {
	return http.StatusNotFound
}

type LoanApplicationFailedError struct{}

func (e LoanApplicationFailedError) Error() string {
	return "loan application failed"
}

func (e LoanApplicationFailedError) Code() int {
	return http.StatusInternalServerError
}

type LoanApprovalFailedError struct{}

func (e LoanApprovalFailedError) Error() string {
	return "loan approval failed"
}

func (e LoanApprovalFailedError) Code() int {
	return http.StatusInternalServerError
}

type LoanPaymentFailedError struct{}

func (e LoanPaymentFailedError) Error() string {
	return "loan payment failed"
}

func (e LoanPaymentFailedError) Code() int {
	return http.StatusInternalServerError
}

type LoanNotPendingError struct{}

func (e LoanNotPendingError) Error() string {
	return "loan is not pending"
}

func (e LoanNotPendingError) Code() int {
	return http.StatusBadRequest
}

type LoanNotApprovedError struct{}

func (e LoanNotApprovedError) Error() string {
	return "loan is not approved"
}

func (e LoanNotApprovedError) Code() int {
	return http.StatusBadRequest
}

type LoanActionNotAllowedError struct{}

func (e LoanActionNotAllowedError) Error() string {
	return "loan action not allowed"
}

func (e LoanActionNotAllowedError) Code() int {
	return http.StatusNotAcceptable
}

type LoanStatusUpdateFailedError struct{}

func (e LoanStatusUpdateFailedError) Error() string {
	return "loan status update failed"
}

func (e LoanStatusUpdateFailedError) Code() int {
	return http.StatusInternalServerError
}
