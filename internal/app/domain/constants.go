package domain

type Role string

const (
	RoleCustomer Role = "Customer"
	RoleApprover Role = "Approver"
)

type Status string

const (
	// Loan Statuses
	LoanStatusPending  Status = "PENDING"
	LoanStatusApproved Status = "APPROVED"
	LoanStatusRejected Status = "REJECTED"
	LoanStatusPaid     Status = "PAID"

	// Payment Statuses
	PayStatusPending Status = "PENDING"
	PayStatusPaid    Status = "PAID"
)
