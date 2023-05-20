package usecase

import (
	"time"

	"github.com/mainakchhari/mini-lender/internal/app/application/errors"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"github.com/mainakchhari/mini-lender/internal/app/domain/repository"
)

type ActionLoanUri struct {
	LoanID int `uri:"id" binding:"required"`
}

type ActionLoanBody struct {
	Action string `form:"action" json:"action" binding:"required"`
}

type ActionLoanArgs struct {
	ActionLoanUri
	ActionLoanBody
	CustomerRepository repository.IUser
	LoanRepository     repository.ILoan
	PaymentRepository  repository.IPayment
}

func ActionLoan(args ActionLoanArgs) errors.BaseError {
	loan, err := args.LoanRepository.Get(args.LoanID)
	if err != nil {
		return &errors.LoanNotFoundError{}
	}
	if loan.Status != domain.LoanStatusPending {
		return &errors.LoanNotPendingError{}
	}
	// check if requested action is either `approved` or `rejected`
	actionAllowed := false
	for _, allowed_status := range []domain.Status{domain.LoanStatusRejected, domain.LoanStatusApproved} {
		if args.Action == string(allowed_status) {
			actionAllowed = true
		}
	}
	if !actionAllowed {
		return &errors.LoanActionNotAllowedError{}
	}
	loan.Status = domain.Status(args.Action)
	loan.ApprovedDate = time.Now()
	_, err = args.LoanRepository.Save(loan)
	if err != nil {
		return &errors.LoanNotApprovedError{}
	}
	return nil
}
