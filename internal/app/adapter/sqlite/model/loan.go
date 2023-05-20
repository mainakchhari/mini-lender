package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Loan model
type Loan struct {
	ID           int             `gorm:"column:loan_id;primaryKey"`
	CustomerID   int             `gorm:"column:customer_id;notnull;index"`
	Customer     User            `gorm:"foreignKey:CustomerID;references:ID"`
	Amount       decimal.Decimal `gorm:"column:amount;type:decimal(15,2);notnull"`
	Status       string          `gorm:"column:status;type:varchar(20);notnull;index"`
	CreatedDate  time.Time       `gorm:"column:created_date;notnull;autoCreateTime;index:,sort:desc"`
	ApproverID   int             `gorm:"column:approver_id;"`
	Approver     *User           `gorm:"foreignKey:ApproverID;references:ID"`
	ApprovedDate time.Time       `gorm:"column:approved_date;"`
}
