package repository

import "github.com/mainakchhari/mini-lender/internal/app/domain"

type IUser interface {
	Get(id int) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	Save(user domain.User) (domain.User, error)
	Authenticate(username string, password string) (domain.User, error)
}
