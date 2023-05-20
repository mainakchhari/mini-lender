package repository

import (
	"time"

	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite/model"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"gorm.io/gorm"
)

type Loan struct{}

func (l Loan) Get(id int) (domain.Loan, error) {
	db := sqlite.Connection()
	var loan model.Loan
	result := db.First(&loan, id)
	if result.Error != nil {
		return domain.Loan{}, result.Error
	}
	return domain.Loan{
		ID:           loan.ID,
		CustomerID:   loan.CustomerID,
		Amount:       loan.Amount,
		Status:       domain.Status(loan.Status),
		CreatedDate:  loan.CreatedDate,
		ApproverID:   loan.ApproverID,
		ApprovedDate: loan.ApprovedDate,
	}, nil
}

func (l Loan) Save(loanEnt domain.Loan) (domain.Loan, error) {
	db := sqlite.Connection()

	loan := model.Loan{
		ID:           loanEnt.ID,
		CustomerID:   loanEnt.CustomerID,
		Amount:       loanEnt.Amount,
		Status:       string(loanEnt.Status),
		ApproverID:   loanEnt.ApproverID,
		ApprovedDate: loanEnt.ApprovedDate,
	}

	var tx *gorm.DB
	if loanEnt.ID == 0 {
		loan.CreatedDate = time.Now()
		tx = db.Omit("ApproverID", "ApprovedDate").Create(&loan)
		loanEnt.ID = loan.ID
		loanEnt.CreatedDate = loan.CreatedDate
	} else {
		loan.ApproverID = loanEnt.ApproverID
		loan.ApprovedDate = loanEnt.ApprovedDate
		tx = db.Omit("CreatedDate").Save(&loan)
	}
	if tx.Error != nil {
		return domain.Loan{}, tx.Error
	}
	return loanEnt, nil
}

func (l Loan) ListFilter(query interface{}, args ...interface{}) ([]domain.Loan, error) {
	db := sqlite.Connection()
	var loans []model.Loan
	result := db.Where(query, args...).Find(&loans)
	if result.Error != nil {
		return []domain.Loan{}, result.Error
	}
	var loanEnts []domain.Loan
	for _, loan := range loans {
		loanEnts = append(loanEnts, domain.Loan{
			ID:           loan.ID,
			CustomerID:   loan.CustomerID,
			Amount:       loan.Amount,
			Status:       domain.Status(loan.Status),
			CreatedDate:  loan.CreatedDate,
			ApproverID:   loan.ApproverID,
			ApprovedDate: loan.ApprovedDate,
		})
	}
	return loanEnts, nil
}
