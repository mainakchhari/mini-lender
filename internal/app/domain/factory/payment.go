package factory

import (
	"time"

	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"github.com/shopspring/decimal"
)

// Payment is a factory of domain.Payment
type Payment struct{}

func (pf Payment) GeneratePaymentsFromLoan(loanEnt domain.Loan, NumInstalments int) []domain.Payment {
	if loanEnt.ID == 0 {
		return []domain.Payment{}
	}
	payments := []domain.Payment{}
	genDate := time.Now()
	for i := 0; i < NumInstalments; i++ {
		payment := domain.Payment{
			ID:          0,
			LoanID:      loanEnt.ID,
			DueAmount:   loanEnt.Amount.Div(decimal.NewFromInt(int64(NumInstalments))),
			DueDate:     genDate.AddDate(0, 0, 7*(i+1)),
			Status:      "PENDING",
			CreatedDate: time.Now(),
		}
		payments = append(payments, payment)
	}
	return payments
}
