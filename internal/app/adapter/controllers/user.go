package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/middleware"
	"github.com/mainakchhari/mini-lender/internal/app/application/usecase"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"github.com/mainakchhari/mini-lender/internal/app/domain/factory"
	"github.com/mainakchhari/mini-lender/internal/app/domain/repository"
)

type UserController struct {
	userRepository    repository.IUser
	loanRepository    repository.ILoan
	paymentRepository repository.IPayment
	paymentFactory    factory.IPayment
}

func SetupUserRoutes(
	r *gin.RouterGroup,
	userRepository repository.IUser,
	loanRepository repository.ILoan,
	paymentRepository repository.IPayment,
	paymentFactory factory.IPayment,
) {
	ctrl := UserController{
		userRepository:    userRepository,
		loanRepository:    loanRepository,
		paymentRepository: paymentRepository,
		paymentFactory:    paymentFactory,
	}
	// Public Route
	r.POST(
		"/register",
		ctrl.createUser,
	)
	r.GET(
		"/:uid/loans",
		middleware.BasicAuthMiddleware(
			[]domain.Role{
				domain.RoleCustomer,
			},
			userRepository,
		),
		ctrl.listLoans,
	)
	r.POST(
		"/:uid/loans",
		middleware.BasicAuthMiddleware(
			[]domain.Role{
				domain.RoleCustomer,
			},
			userRepository,
		),
		ctrl.applyLoan,
	)
	r.PUT(
		"/:uid/loans/:lid/pay",
		middleware.BasicAuthMiddleware(
			[]domain.Role{
				domain.RoleCustomer,
			},
			userRepository,
		),
		ctrl.makePayment,
	)
}

func (ctrl UserController) createUser(c *gin.Context) {
	body := usecase.CreateUserBody{}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	args := usecase.CreateUserArgs{
		CreateUserBody: body,
		UserRepository: ctrl.userRepository,
	}
	order, err := usecase.CreateUser(args)
	if err != nil {
		c.AbortWithStatusJSON(err.Code(), gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, order)
}

func (ctrl UserController) listLoans(c *gin.Context) {
	var uri usecase.ListLoansUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	args := usecase.ListLoansArgs{
		ListLoansUri:   uri,
		LoanRepository: ctrl.loanRepository,
	}
	order, err := usecase.ListLoans(args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, order)
}

func (ctrl UserController) applyLoan(c *gin.Context) {
	var uri usecase.ApplyLoanUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	body := usecase.ApplyLoanBody{}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	args := usecase.ApplyLoanArgs{
		ApplyLoanBody:      body,
		CustomerRepository: ctrl.userRepository,
		LoanRepository:     ctrl.loanRepository,
		PaymentRepository:  ctrl.paymentRepository,
		PaymentFactory:     ctrl.paymentFactory,
	}
	order, err := usecase.ApplyLoan(args)
	if err != nil {
		c.AbortWithStatusJSON(err.Code(), gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, order)
}

func (ctrl UserController) makePayment(c *gin.Context) {
	var uri usecase.MakePaymentUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	body := usecase.MakePaymentBody{}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	args := usecase.MakePaymentArgs{
		MakePaymentUri:    uri,
		MakePaymentBody:   body,
		LoanRepository:    ctrl.loanRepository,
		PaymentRepository: ctrl.paymentRepository,
	}
	err := usecase.MakePayment(args)
	if err != nil {
		c.AbortWithStatusJSON(err.Code(), gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
