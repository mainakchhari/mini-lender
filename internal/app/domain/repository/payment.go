package repository

import "github.com/mainakchhari/mini-lender/internal/app/domain"

type IPayment interface {
	Get(id int) (domain.Payment, error)
	Save(payment domain.Payment) (domain.Payment, error)
	GetNextByLoanId(loanId int) (domain.Payment, error)
}
