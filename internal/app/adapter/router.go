package adapter

import (
	"github.com/gin-gonic/gin"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/controllers"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/repository"
	"github.com/mainakchhari/mini-lender/internal/app/domain/factory"
)

var (
	userRepository    = repository.User{}
	loanRepository    = repository.Loan{}
	paymentRepository = repository.Payment{}
	paymentFactory    = factory.Payment{}
)

func Router() *gin.Engine {
	r := gin.Default()
	userGroup := r.Group("/users")
	loanGroup := r.Group("/loans")
	controllers.SetupUserRoutes(userGroup, userRepository, loanRepository, paymentRepository, paymentFactory)
	controllers.SetupLoanRoutes(loanGroup, userRepository, loanRepository)
	r.GET("/", func(c *gin.Context) { c.String(200, "OK") })
	return r
}
