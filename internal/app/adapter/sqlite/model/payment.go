package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Payment model
type Payment struct {
	ID          int             `gorm:"column:payment_id;primaryKey"`
	LoanID      int             `gorm:"column:loan_id;notnull;index"`
	Loan        Loan            `gorm:"foreignKey:LoanID;references:ID"`
	DueAmount   decimal.Decimal `gorm:"column:due_amount;type:decimal(15,2);notnull"`
	DueDate     time.Time       `gorm:"column:due_date;notnull"`
	PaidAmount  decimal.Decimal `gorm:"column:paid_amount;type:decimal(15,2)"`
	PaidDate    time.Time       `gorm:"column:paid_date;"`
	Status      string          `gorm:"column:status;type:varchar(20);notnull;index"`
	CreatedDate time.Time       `gorm:"column:created_date;notnull;autoCreateTime;index:,sort:desc"`
}
