package usecase

import (
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/mainakchhari/mini-lender/internal/app/application/errors"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	mockFactory "github.com/mainakchhari/mini-lender/internal/app/domain/mock_factory"
	mockRepo "github.com/mainakchhari/mini-lender/internal/app/domain/mock_repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type listLoansTest struct {
	id       string
	args     ListLoansArgs
	expected []LoanDetailResponse
	err      error
	preTest  func()
}

func TestListLoans(t *testing.T) {
	ctrl := gomock.NewController(t)
	loanMock := mockRepo.NewMockILoan(ctrl)

	tests := []listLoansTest{
		{
			id: "fails when user not found",
			args: ListLoansArgs{
				ListLoansUri: ListLoansUri{
					CustomerID: 1,
				},
				LoanRepository: loanMock,
			},
			expected: []LoanDetailResponse{},
			err:      gorm.ErrRecordNotFound,
			preTest: func() {
				loanMock.EXPECT().ListFilter("customer_id = ?", 1).Return([]domain.Loan{}, gorm.ErrRecordNotFound)
			},
		},
		{
			id: "success when loans are found",
			args: ListLoansArgs{
				ListLoansUri: ListLoansUri{
					CustomerID: 1,
				},
				LoanRepository: loanMock,
			},
			expected: []LoanDetailResponse{
				{
					ID:          1,
					CustomerID:  1,
					Amount:      decimal.NewFromInt(1000),
					Status:      "PENDING",
					CreatedDate: time.Now(),
				},
			},
			err: nil,
			preTest: func() {
				//monkey patch (replace) time.now with a function that returns a fixed time
				monkey.Patch(time.Now, func() time.Time {
					return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				})

				domainLoans := []domain.Loan{
					{
						ID:          1,
						CustomerID:  1,
						Amount:      decimal.NewFromInt(1000),
						Status:      domain.LoanStatusPending,
						CreatedDate: time.Now(),
					},
				}
				loanMock.EXPECT().ListFilter("customer_id = ?", 1).Return(domainLoans, nil)
			},
		},
		{
			id: "success when no loans are found",
			args: ListLoansArgs{
				ListLoansUri: ListLoansUri{
					CustomerID: 1,
				},
				LoanRepository: loanMock,
			},
			expected: []LoanDetailResponse{},
			err:      nil,
			preTest: func() {
				loanMock.EXPECT().ListFilter("customer_id = ?", 1).Return([]domain.Loan{}, nil)
			},
		},
	}

	for _, test := range tests {
		test.preTest()
		if output, err := ListLoans(test.args); !reflect.DeepEqual(test.expected, output) || err != test.err {
			t.Errorf("%s:\nArgs: %+v\nOutput: %+v\nExpected: %+v\n", test.id, test.args, output, test.expected)
		}
	}

}

type applyLoanTest struct {
	id       string
	args     ApplyLoanArgs
	expected ApplyLoanResponse
	err      errors.BaseError
	preTest  func()
}

func TestApplyLoan(t *testing.T) {
	ctrl := gomock.NewController(t)
	userMock := mockRepo.NewMockIUser(ctrl)
	loanMock := mockRepo.NewMockILoan(ctrl)
	paymentMock := mockRepo.NewMockIPayment(ctrl)
	payFactoryMock := mockFactory.NewMockIPayment(ctrl)

	tests := []applyLoanTest{
		{
			id: "fails when user not found",
			args: ApplyLoanArgs{
				ApplyLoanUri: ApplyLoanUri{
					CustomerID: 1,
				},
				ApplyLoanBody:      ApplyLoanBody{},
				CustomerRepository: userMock,
			},
			expected: ApplyLoanResponse{},
			err:      &errors.UserNotFoundError{},
			preTest: func() {
				userMock.EXPECT().Get(1).Return(domain.User{}, gorm.ErrRecordNotFound)
			},
		},
		{
			id: "fails when loan save fails",
			args: ApplyLoanArgs{
				ApplyLoanUri: ApplyLoanUri{
					CustomerID: 1,
				},
				ApplyLoanBody: ApplyLoanBody{
					LoanAmount:     decimal.NewFromInt(1000),
					NumInstalments: 10,
				},
				CustomerRepository: userMock,
				LoanRepository:     loanMock,
			},
			expected: ApplyLoanResponse{},
			err:      &errors.LoanApplicationFailedError{},
			preTest: func() {
				domainLoan := domain.Loan{CustomerID: 1, Amount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				userMock.EXPECT().Get(1).Return(domain.User{ID: 1}, nil)
				loanMock.EXPECT().Save(domainLoan).Return(domain.Loan{}, gorm.ErrInvalidData)
			},
		},
		{
			id: "fails when payment save fails",
			args: ApplyLoanArgs{
				ApplyLoanUri: ApplyLoanUri{
					CustomerID: 1,
				},
				ApplyLoanBody: ApplyLoanBody{
					LoanAmount:     decimal.NewFromInt(1000),
					NumInstalments: 1,
				},
				CustomerRepository: userMock,
				LoanRepository:     loanMock,
				PaymentRepository:  paymentMock,
				PaymentFactory:     payFactoryMock,
			},
			expected: ApplyLoanResponse{},
			err:      &errors.LoanApplicationFailedError{},
			preTest: func() {
				//monkey patch (replace) time.now with a function that returns a fixed time
				monkey.Patch(time.Now, func() time.Time {
					return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				})

				userMock.EXPECT().Get(1).Return(domain.User{ID: 1}, nil)
				domainLoan := domain.Loan{CustomerID: 1, Amount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				createdLoan := domainLoan
				createdLoan.ID = 1
				loanMock.EXPECT().Save(domainLoan).Return(createdLoan, nil)
				domainPayment := domain.Payment{
					ID:          1,
					LoanID:      1,
					DueAmount:   decimal.NewFromInt(1000),
					DueDate:     time.Now().AddDate(0, 0, 7),
					Status:      domain.PayStatusPending,
					CreatedDate: time.Now(),
				}
				payFactoryMock.EXPECT().GeneratePaymentsFromLoan(createdLoan, 1).Return([]domain.Payment{domainPayment})
				paymentMock.EXPECT().Save(domainPayment).Return(domain.Payment{}, gorm.ErrInvalidData)
			},
		},
		{
			id: "success when loan application is successful",
			args: ApplyLoanArgs{
				ApplyLoanUri: ApplyLoanUri{
					CustomerID: 1,
				},
				ApplyLoanBody: ApplyLoanBody{
					LoanAmount:     decimal.NewFromInt(1000),
					NumInstalments: 1,
				},
				CustomerRepository: userMock,
				LoanRepository:     loanMock,
				PaymentRepository:  paymentMock,
				PaymentFactory:     payFactoryMock,
			},
			expected: ApplyLoanResponse{
				LoanID: 1,
				ApplyLoanBody: ApplyLoanBody{
					LoanAmount:     decimal.NewFromInt(1000),
					NumInstalments: 1,
				},
				Status:      string(domain.LoanStatusPending),
				CreatedDate: time.Now(),
			},
			err: nil,
			preTest: func() {
				//monkey patch (replace) time.now with a function that returns a fixed time
				monkey.Patch(time.Now, func() time.Time {
					return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				})

				userMock.EXPECT().Get(1).Return(domain.User{ID: 1}, nil)
				domainLoan := domain.Loan{CustomerID: 1, Amount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				createdLoan := domainLoan
				createdLoan.ID = 1
				createdLoan.CreatedDate = time.Now()
				loanMock.EXPECT().Save(domainLoan).Return(createdLoan, nil)
				domainPayment := domain.Payment{
					ID:          1,
					LoanID:      1,
					DueAmount:   decimal.NewFromInt(1000),
					DueDate:     time.Now().AddDate(0, 0, 7),
					Status:      domain.PayStatusPending,
					CreatedDate: time.Now(),
				}
				createdPayment := domainPayment
				createdPayment.ID = 1
				createdPayment.CreatedDate = time.Now()
				payFactoryMock.EXPECT().GeneratePaymentsFromLoan(createdLoan, 1).Return([]domain.Payment{domainPayment})
				paymentMock.EXPECT().Save(domainPayment).Return(createdPayment, nil)
			},
		},
	}

	for _, test := range tests {
		test.preTest()
		if output, err := ApplyLoan(test.args); !reflect.DeepEqual(test.expected, output) || err != test.err {
			t.Errorf("%s:\nArgs: %+v\nOutput: %+v\nExpected: %+v\n", test.id, test.args, output, test.expected)
		}
	}
}

type makePaymentTest struct {
	id      string
	args    MakePaymentArgs
	err     errors.BaseError
	preTest func()
}

func TestMakePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	loanMock := mockRepo.NewMockILoan(ctrl)
	paymentMock := mockRepo.NewMockIPayment(ctrl)

	tests := []makePaymentTest{
		{
			id: "fails when loan not found for customer id",
			args: MakePaymentArgs{
				MakePaymentUri: MakePaymentUri{
					UserID: 1,
					LoanID: 1,
				},
				MakePaymentBody: MakePaymentBody{},
				LoanRepository:  loanMock,
			},
			err: errors.LoanNotFoundError{},
			preTest: func() {
				loanMock.EXPECT().ListFilter("customer_id = ? AND loan_id = ?", 1, 1).Return([]domain.Loan{}, nil)
			},
		},
		{
			id: "fails when no pending payment for loan id",
			args: MakePaymentArgs{
				MakePaymentUri: MakePaymentUri{
					UserID: 1,
					LoanID: 1,
				},
				MakePaymentBody:   MakePaymentBody{},
				LoanRepository:    loanMock,
				PaymentRepository: paymentMock,
			},
			err: errors.PaymentNotFoundError{},
			preTest: func() {
				domainLoan := domain.Loan{ID: 1, CustomerID: 1, Amount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				loanMock.EXPECT().ListFilter("customer_id = ? AND loan_id = ?", 1, 1).Return([]domain.Loan{domainLoan}, nil)
				paymentMock.EXPECT().GetNextByLoanId(1).Return(domain.Payment{}, gorm.ErrRecordNotFound)
			},
		},
		{
			id: "fails when payment amount less than due amount",
			args: MakePaymentArgs{
				MakePaymentUri: MakePaymentUri{
					UserID: 1,
					LoanID: 1,
				},
				MakePaymentBody: MakePaymentBody{
					Amount: decimal.NewFromInt(500),
				},
				LoanRepository:    loanMock,
				PaymentRepository: paymentMock,
			},
			err: errors.InvalidPaymentAmountError{},
			preTest: func() {
				domainLoan := domain.Loan{ID: 1, CustomerID: 1, Amount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				domainPayment := domain.Payment{ID: 1, LoanID: 1, DueAmount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				loanMock.EXPECT().ListFilter("customer_id = ? AND loan_id = ?", 1, 1).Return([]domain.Loan{domainLoan}, nil)
				paymentMock.EXPECT().GetNextByLoanId(1).Return(domainPayment, nil)
			},
		},
		{
			id: "fails when payment db save fails",
			args: MakePaymentArgs{
				MakePaymentUri: MakePaymentUri{
					UserID: 1,
					LoanID: 1,
				},
				MakePaymentBody: MakePaymentBody{
					Amount: decimal.NewFromInt(1000),
				},
				LoanRepository:    loanMock,
				PaymentRepository: paymentMock,
			},
			err: errors.LoanPaymentFailedError{},
			preTest: func() {
				domainLoan := domain.Loan{ID: 1, CustomerID: 1, Amount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				domainPayment := domain.Payment{ID: 1, LoanID: 1, DueAmount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending, CreatedDate: time.Now()}
				loanMock.EXPECT().ListFilter("customer_id = ? AND loan_id = ?", 1, 1).Return([]domain.Loan{domainLoan}, nil)
				paymentMock.EXPECT().GetNextByLoanId(1).Return(domainPayment, nil)
				savedPayment := domainPayment
				savedPayment.PaidAmount = decimal.NewFromInt(1000)
				savedPayment.Status = domain.PayStatusPaid
				savedPayment.PaidDate = time.Now()
				paymentMock.EXPECT().Save(savedPayment).Return(domain.Payment{}, gorm.ErrInvalidData)
			},
		},
		{
			id: "success when loan repayment is successful",
			args: MakePaymentArgs{
				MakePaymentUri: MakePaymentUri{
					UserID: 1,
					LoanID: 1,
				},
				MakePaymentBody: MakePaymentBody{
					Amount: decimal.NewFromInt(1000),
				},
				LoanRepository:    loanMock,
				PaymentRepository: paymentMock,
			},
			err: nil,
			preTest: func() {
				domainLoan := domain.Loan{ID: 1, CustomerID: 1, Amount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending}
				domainPayment := domain.Payment{ID: 1, LoanID: 1, DueAmount: decimal.NewFromInt(1000), Status: domain.LoanStatusPending, CreatedDate: time.Now()}
				loanMock.EXPECT().ListFilter("customer_id = ? AND loan_id = ?", 1, 1).Return([]domain.Loan{domainLoan}, nil)
				paymentMock.EXPECT().GetNextByLoanId(1).Return(domainPayment, nil)
				savedPayment := domainPayment
				savedPayment.PaidAmount = decimal.NewFromInt(1000)
				savedPayment.Status = domain.PayStatusPaid
				savedPayment.PaidDate = time.Now()
				paymentMock.EXPECT().Save(savedPayment).Return(savedPayment, nil)
				paymentMock.EXPECT().GetNextByLoanId(1).Return(domain.Payment{}, gorm.ErrRecordNotFound)
				domainLoan.Status = domain.LoanStatusPaid
				loanMock.EXPECT().Save(domainLoan).Return(domainLoan, nil)
			},
		},
	}

	for _, test := range tests {
		test.preTest()
		if err := MakePayment(test.args); err != test.err {
			t.Errorf("%s:\nArgs: %+v\nOutput: %+v\nExpected: %+v\n", test.id, test.args, err, test.err)
		}
	}

}
