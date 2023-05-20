package adapter

import (
	"github.com/gin-gonic/gin"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/controllers"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/repository"
)

var (
	userRepository    = repository.User{}
	loanRepository    = repository.Loan{}
	paymentRepository = repository.Payment{}
)

func Router() *gin.Engine {
	r := gin.Default()
	userGroup := r.Group("/users")
	loanGroup := r.Group("/loans")
	controllers.SetupUserRoutes(userGroup, userRepository, loanRepository, paymentRepository)
	controllers.SetupLoanRoutes(loanGroup, userRepository, loanRepository, paymentRepository)
	r.GET("/", func(c *gin.Context) { c.String(200, "OK") })
	return r
}
