package repository

import "github.com/mainakchhari/mini-lender/internal/app/domain"

type ILoan interface {
	Get(id int) (domain.Loan, error)
	Save(loan domain.Loan) (domain.Loan, error)
	ListFilter(query interface{}, args ...interface{}) ([]domain.Loan, error)
}
