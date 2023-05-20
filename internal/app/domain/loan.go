package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Loan Entity
type Loan struct {
	ID           int
	CustomerID   int
	Customer     User
	Amount       decimal.Decimal
	Status       Status
	CreatedDate  time.Time
	ApproverID   int
	Approver     User
	ApprovedDate time.Time
}
