package repository

import (
	"fmt"

	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite"
	"github.com/mainakchhari/mini-lender/internal/app/adapter/sqlite/model"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct{}

func (u User) Get(id int) (domain.User, error) {
	db := sqlite.Connection()
	var user model.User
	result := db.First(&user, id)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return domain.User{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}, nil
}

func (u User) GetByUsername(username string) (domain.User, error) {
	db := sqlite.Connection()
	var user model.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return domain.User{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}, nil
}

func (u User) Save(userEnt domain.User) (domain.User, error) {
	db := sqlite.Connection()

	user := model.User{
		ID:       userEnt.ID,
		Username: userEnt.Username,
		Name:     userEnt.Name,
		Role:     userEnt.Role,
		Password: userEnt.Password,
	}

	// allows only valid role values
	var roleAllowed bool = false
	for _, role := range []domain.Role{domain.RoleCustomer, domain.RoleApprover} {
		if userEnt.Role == string(role) {
			roleAllowed = true
		}
	}
	if !roleAllowed {
		return domain.User{}, fmt.Errorf("invalid role value: %w", gorm.ErrInvalidData)
	}

	var tx *gorm.DB
	if userEnt.ID == 0 {
		tx = db.Create(&user)
		userEnt.ID = user.ID
	} else {
		tx = db.Save(&user)
	}
	if tx.Error != nil {
		return domain.User{}, tx.Error
	}
	return userEnt, nil
}

func (u User) Authenticate(username string, password string) (domain.User, error) {
	db := sqlite.Connection()
	var user model.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}, nil
}
