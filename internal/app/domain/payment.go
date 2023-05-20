package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Payment Entity
type Payment struct {
	ID          int
	LoanID      int
	Loan        Loan
	DueAmount   decimal.Decimal
	DueDate     time.Time
	PaidAmount  decimal.Decimal
	PaidDate    time.Time
	Status      Status
	CreatedDate time.Time
}
