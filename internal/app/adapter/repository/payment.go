package repository

import (
	"time"

	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite/model"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"gorm.io/gorm"
)

type Payment struct{}

func (p Payment) Get(id int) (domain.Payment, error) {
	db := sqlite.Connection()
	var payment model.Payment
	result := db.First(&payment, id)
	if result.Error != nil {
		return domain.Payment{}, result.Error
	}
	return domain.Payment{
		ID:          payment.ID,
		LoanID:      payment.LoanID,
		DueAmount:   payment.DueAmount,
		DueDate:     payment.DueDate,
		PaidAmount:  payment.PaidAmount,
		PaidDate:    payment.PaidDate,
		Status:      domain.Status(payment.Status),
		CreatedDate: payment.CreatedDate,
	}, nil
}

func (p Payment) Save(paymentEnt domain.Payment) (domain.Payment, error) {
	db := sqlite.Connection()

	payment := model.Payment{
		ID:         paymentEnt.ID,
		LoanID:     paymentEnt.LoanID,
		DueAmount:  paymentEnt.DueAmount,
		DueDate:    paymentEnt.DueDate,
		Status:     string(paymentEnt.Status),
		PaidAmount: paymentEnt.PaidAmount,
		PaidDate:   paymentEnt.PaidDate,
	}

	var tx *gorm.DB
	if paymentEnt.ID == 0 {
		payment.CreatedDate = time.Now()
		tx = db.Omit("PaidAmount", "PaidDate").Create(&payment)
		paymentEnt.ID = payment.ID
		paymentEnt.CreatedDate = payment.CreatedDate
	} else {
		tx = db.Omit("CreatedDate").Save(&payment)
	}
	if tx.Error != nil {
		return domain.Payment{}, tx.Error
	}
	return paymentEnt, nil
}

func (p Payment) GetNextByLoanId(loanId int) (domain.Payment, error) {
	db := sqlite.Connection()
	var payment model.Payment
	result := db.Where("loan_id = ? AND status = ?", loanId, "PENDING").Order("due_date asc").First(&payment)
	if result.Error != nil {
		return domain.Payment{}, result.Error
	}
	return domain.Payment{
		ID:          payment.ID,
		LoanID:      payment.LoanID,
		DueAmount:   payment.DueAmount,
		DueDate:     payment.DueDate,
		PaidAmount:  payment.PaidAmount,
		PaidDate:    payment.PaidDate,
		Status:      domain.Status(payment.Status),
		CreatedDate: payment.CreatedDate,
	}, nil
}
