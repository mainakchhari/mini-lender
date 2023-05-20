package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/middleware"
	"github.com/mainakchhari/mini-lender/internal/app/application/usecase"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"github.com/mainakchhari/mini-lender/internal/app/domain/repository"
)

type LoanController struct {
	loanRepository repository.ILoan
}

func SetupLoanRoutes(
	r *gin.RouterGroup,
	userRepository repository.IUser,
	loanRepository repository.ILoan,
) {
	ctrl := LoanController{
		loanRepository: loanRepository,
	}
	r.PUT(
		"/:id/action",
		middleware.BasicAuthMiddleware(
			[]domain.Role{
				domain.RoleApprover,
			},
			userRepository,
		),
		ctrl.actionLoan,
	)
}

func (ctrl LoanController) actionLoan(c *gin.Context) {
	var uri usecase.ActionLoanUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	body := usecase.ActionLoanBody{}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Validation Error"})
		return
	}
	args := usecase.ActionLoanArgs{
		ActionLoanUri:  uri,
		ActionLoanBody: body,
		LoanRepository: ctrl.loanRepository,
	}
	err := usecase.ActionLoan(args)
	if err != nil {
		c.AbortWithStatusJSON(err.Code(), err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
