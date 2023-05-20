package usecase

import (
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/mainakchhari/mini-lender/internal/app/application/errors"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	mockRepo "github.com/mainakchhari/mini-lender/internal/app/domain/mock_repository"
	"gorm.io/gorm"
)

type actionLoanTest struct {
	id       string
	args     ActionLoanArgs
	expected errors.BaseError
	preTest  func()
}

func TestActionLoan(t *testing.T) {
	ctrl := gomock.NewController(t)
	loanMock := mockRepo.NewMockILoan(ctrl)

	tests := []actionLoanTest{
		{
			id: "fails when loan is not found",
			args: ActionLoanArgs{
				ActionLoanUri: ActionLoanUri{
					LoanID: 1,
				},
				ActionLoanBody: ActionLoanBody{
					Action: "APPROVED",
				},
				LoanRepository: loanMock,
			},
			expected: &errors.LoanNotFoundError{},
			preTest: func() {
				loanMock.EXPECT().Get(1).Return(domain.Loan{}, gorm.ErrRecordNotFound)
			},
		},
		{
			id: "fails when loan is not pending",
			args: ActionLoanArgs{
				ActionLoanUri: ActionLoanUri{
					LoanID: 1,
				},
				ActionLoanBody: ActionLoanBody{
					Action: "APPROVED",
				},
				LoanRepository: loanMock,
			},
			expected: &errors.LoanNotPendingError{},
			preTest: func() {
				loanMock.EXPECT().Get(1).Return(domain.Loan{Status: domain.LoanStatusApproved}, nil)
			},
		},
		{
			id: "fails when loan action is not allowed/recognised",
			args: ActionLoanArgs{
				ActionLoanUri: ActionLoanUri{
					LoanID: 1,
				},
				ActionLoanBody: ActionLoanBody{
					Action: "randomstring",
				},
				LoanRepository: loanMock,
			},
			expected: &errors.LoanActionNotAllowedError{},
			preTest: func() {
				loanMock.EXPECT().Get(1).Return(domain.Loan{Status: domain.LoanStatusPending}, nil)
			},
		},
		{
			id: "fails when db save fails",
			args: ActionLoanArgs{
				ActionLoanUri: ActionLoanUri{
					LoanID: 1,
				},
				ActionLoanBody: ActionLoanBody{
					Action: "REJECTED",
				},
				LoanRepository: loanMock,
			},
			expected: &errors.LoanNotApprovedError{},
			preTest: func() {
				//monkey patch (replace) time.now with a function that returns a fixed time
				monkey.Patch(time.Now, func() time.Time {
					return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				})

				domainLoan := domain.Loan{ID: 1, Status: domain.LoanStatusPending}
				loanMock.EXPECT().Get(1).Return(domainLoan, nil)
				domainLoan.Status = domain.LoanStatusRejected
				domainLoan.ApprovedDate = time.Now()
				loanMock.EXPECT().Save(domainLoan).Return(domainLoan, gorm.ErrRecordNotFound)
			},
		},
		{
			id: "success when loan is approved",
			args: ActionLoanArgs{
				ActionLoanUri: ActionLoanUri{
					LoanID: 1,
				},
				ActionLoanBody: ActionLoanBody{
					Action: "APPROVED",
				},
				LoanRepository: loanMock,
			},
			expected: nil,
			preTest: func() {
				//monkey patch (replace) time.now with a function that returns a fixed time
				monkey.Patch(time.Now, func() time.Time {
					return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				})

				domainLoan := domain.Loan{ID: 1, Status: domain.LoanStatusPending}
				loanMock.EXPECT().Get(1).Return(domainLoan, nil)
				domainLoan.Status = domain.LoanStatusApproved
				domainLoan.ApprovedDate = time.Now()
				loanMock.EXPECT().Save(domainLoan).Return(domainLoan, nil)
			},
		},
		{
			id: "success when loan is rejected",
			args: ActionLoanArgs{
				ActionLoanUri: ActionLoanUri{
					LoanID: 1,
				},
				ActionLoanBody: ActionLoanBody{
					Action: "REJECTED",
				},
				LoanRepository: loanMock,
			},
			expected: nil,
			preTest: func() {
				//monkey patch (replace) time.now with a function that returns a fixed time
				monkey.Patch(time.Now, func() time.Time {
					return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				})

				domainLoan := domain.Loan{ID: 1, Status: domain.LoanStatusPending}
				loanMock.EXPECT().Get(1).Return(domainLoan, nil)
				domainLoan.Status = domain.LoanStatusRejected
				domainLoan.ApprovedDate = time.Now()
				loanMock.EXPECT().Save(domainLoan).Return(domainLoan, nil)
			},
		},
	}

	for _, test := range tests {
		test.preTest()
		if output := ActionLoan(test.args); output != test.expected {
			t.Errorf("%s:\nArgs: %+v\nOutput: %+v\nExpected: %+v\n", test.id, test.args, output, test.expected)
		}
	}
}
