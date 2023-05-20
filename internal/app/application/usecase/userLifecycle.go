package usecase

import (
	"github.com/mainakchhari/mini-lender/internal/app/application/errors"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"github.com/mainakchhari/mini-lender/internal/app/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserBody struct {
	Username string `form:"username" json:"username" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Role     string `form:"role" json:"role" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type CreateUserArgs struct {
	CreateUserBody
	UserRepository repository.IUser
}

type CreateUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func CreateUser(args CreateUserArgs) (CreateUserResponse, errors.BaseError) {
	var roleAllowed bool
	for _, allowedRole := range []domain.Role{domain.RoleApprover, domain.RoleCustomer} {
		if string(allowedRole) == args.Role {
			roleAllowed = true
		}
	}
	if !roleAllowed {
		return CreateUserResponse{}, errors.UserRoleInvalidError{}
	}

	// Check if username exists already
	_, err := args.UserRepository.GetByUsername(args.Username)
	if err == nil {
		return CreateUserResponse{}, errors.UsernameAlreadyExistsError{}
	}

	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return CreateUserResponse{}, errors.InvalidPasswordStringError{}
	}

	user := domain.User{
		Username: args.Username,
		Name:     args.Name,
		Role:     args.Role,
		Password: string(hash),
	}
	user, err = args.UserRepository.Save(user)
	if err != nil {
		return CreateUserResponse{}, errors.UserCreationFailed{}
	}
	return CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}, nil
}
