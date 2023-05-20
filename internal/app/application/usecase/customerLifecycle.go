package usecase

import (
	"time"

	"github.com/mainakchhari/mini-lender/internal/app/application/errors"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"github.com/mainakchhari/mini-lender/internal/app/domain/factory"
	"github.com/mainakchhari/mini-lender/internal/app/domain/repository"
	"github.com/shopspring/decimal"
)

type ListLoansUri struct {
	CustomerID int `uri:"uid" binding:"required"`
}

type ListLoansArgs struct {
	ListLoansUri
	LoanRepository repository.ILoan
}

type LoanDetailResponse struct {
	ID           int             `json:"loan_id"`
	CustomerID   int             `json:"customer_id"`
	Amount       decimal.Decimal `json:"amount"`
	Status       string          `json:"status"`
	CreatedDate  time.Time       `json:"created_date"`
	ApproverID   int             `json:"approver_id"`
	ApprovedDate time.Time       `json:"approved_date"`
}

func ListLoans(args ListLoansArgs) ([]LoanDetailResponse, error) {
	// Get filtered list of loans by customer_id
	loans, err := args.LoanRepository.ListFilter("customer_id = ?", args.CustomerID)
	if err != nil {
		return []LoanDetailResponse{}, err
	}
	listLoansResponse := []LoanDetailResponse{}
	for _, loan := range loans {
		listLoansResponse = append(listLoansResponse, LoanDetailResponse{
			ID:           loan.ID,
			CustomerID:   loan.CustomerID,
			Amount:       loan.Amount,
			Status:       string(loan.Status),
			CreatedDate:  loan.CreatedDate,
			ApproverID:   loan.ApproverID,
			ApprovedDate: loan.ApprovedDate,
		})
	}
	return listLoansResponse, nil
}

type ApplyLoanUri struct {
	CustomerID int `uri:"uid" binding:"required"`
}

type ApplyLoanBody struct {
	LoanAmount     decimal.Decimal `form:"loan_amount" json:"loan_amount" binding:"required"`
	NumInstalments int             `form:"num_instalments" json:"num_instalments" binding:"required"`
}

type ApplyLoanArgs struct {
	ApplyLoanUri
	ApplyLoanBody
	CustomerRepository repository.IUser
	LoanRepository     repository.ILoan
	PaymentRepository  repository.IPayment
}

type ApplyLoanResponse struct {
	LoanID int `json:"loan_id"`
	ApplyLoanBody
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
}

func ApplyLoan(args ApplyLoanArgs) (ApplyLoanResponse, errors.BaseError) {
	customer, err := args.CustomerRepository.Get(args.CustomerID)
	if err != nil {
		return ApplyLoanResponse{}, &errors.UserNotFoundError{}
	}
	loan := domain.Loan{
		CustomerID: customer.ID,
		Amount:     args.LoanAmount,
		Status:     domain.LoanStatusPending,
	}
	loan, err = args.LoanRepository.Save(loan)
	if err != nil {
		return ApplyLoanResponse{}, &errors.LoanApplicationFailedError{}
	}
	paymentFactory := factory.Payment{}
	for _, payment := range paymentFactory.GeneratePaymentsFromLoan(loan, args.NumInstalments) {
		_, err := args.PaymentRepository.Save(payment)
		if err != nil {
			return ApplyLoanResponse{}, &errors.LoanApplicationFailedError{}
		}
	}
	loan, err = args.LoanRepository.Get(loan.ID)
	if err != nil {
		return ApplyLoanResponse{}, &errors.LoanNotFoundError{}
	}
	return ApplyLoanResponse{
		LoanID:        loan.ID,
		ApplyLoanBody: args.ApplyLoanBody,
		Status:        string(loan.Status),
		CreatedDate:   loan.CreatedDate,
	}, nil
}

type MakePaymentUri struct {
	UserID int `uri:"uid" binding:"required"`
	LoanID int `uri:"lid" binding:"required"`
}

type MakePaymentBody struct {
	Amount decimal.Decimal `form:"amount" json:"amount" binding:"required"`
}

type MakePaymentArgs struct {
	MakePaymentUri
	MakePaymentBody
	LoanRepository    repository.ILoan
	PaymentRepository repository.IPayment
}

func MakePayment(args MakePaymentArgs) errors.BaseError {
	loans, err := args.LoanRepository.ListFilter("customer_id = ? AND loan_id = ?", args.UserID, args.LoanID)
	if err != nil || len(loans) == 0 {
		return errors.LoanNotFoundError{}
	}
	payment, err := args.PaymentRepository.GetNextByLoanId(loans[0].ID)
	if err != nil {
		return errors.PaymentNotFoundError{}
	}
	if (payment.DueAmount).GreaterThan(args.Amount) {
		return errors.InvalidPaymentAmountError{}
	}
	payment.PaidAmount = args.Amount
	payment.Status = domain.PayStatusPaid
	payment.PaidDate = time.Now()
	_, err = args.PaymentRepository.Save(payment)
	if err != nil {
		return errors.LoanPaymentFailedError{}
	}
	return nil
}
