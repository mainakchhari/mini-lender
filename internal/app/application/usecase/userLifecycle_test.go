package usecase

import (
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/mainakchhari/mini-lender/internal/app/application/errors"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	mockRepo "github.com/mainakchhari/mini-lender/internal/app/domain/mock_repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type createUserTest struct {
	id       string
	args     CreateUserArgs
	expected CreateUserResponse
	err      errors.BaseError
	preTest  func()
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	userMock := mockRepo.NewMockIUser(ctrl)

	tests := []createUserTest{
		{
			id: "fail on invalid/disallowed role",
			args: CreateUserArgs{
				CreateUserBody: CreateUserBody{
					Role: "invalid",
				},
				UserRepository: userMock,
			},
			expected: CreateUserResponse{},
			err:      errors.UserRoleInvalidError{},
			preTest:  func() {},
		},
		{
			id: "fail on username already exists",
			args: CreateUserArgs{
				CreateUserBody: CreateUserBody{
					Username: "testuser",
					Role:     string(domain.RoleApprover),
				},
				UserRepository: userMock,
			},
			expected: CreateUserResponse{},
			err:      errors.UsernameAlreadyExistsError{},
			preTest: func() {
				domainUser := domain.User{ID: 1, Username: "testuser", Role: string(domain.RoleApprover)}
				userMock.EXPECT().GetByUsername("testuser").Return(domainUser, nil)
			},
		},
		{
			id: "fail when user db save fails",
			args: CreateUserArgs{
				CreateUserBody: CreateUserBody{
					Username: "testuser",
					Role:     string(domain.RoleApprover),
					Password: "testpassword",
				},
				UserRepository: userMock,
			},
			expected: CreateUserResponse{},
			err:      errors.UserCreationFailed{},
			preTest: func() {
				//monkey patch (replace) bcrypt.GenerateFromPassword with a function that returns a fixed output
				monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
					return []byte("testhash"), nil
				})

				// Generate "hash" to store from user password
				hash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
				domainUser := domain.User{Username: "testuser", Role: string(domain.RoleApprover), Password: string(hash)}
				userMock.EXPECT().GetByUsername("testuser").Return(domain.User{}, gorm.ErrRecordNotFound)
				userMock.EXPECT().Save(domainUser).Return(domain.User{}, gorm.ErrInvalidData)
			},
		},
		{
			id: "success when user created",
			args: CreateUserArgs{
				CreateUserBody: CreateUserBody{
					Username: "testuser",
					Role:     string(domain.RoleApprover),
					Password: "testpassword",
				},
				UserRepository: userMock,
			},
			expected: CreateUserResponse{
				ID:       1,
				Username: "testuser",
				Role:     string(domain.RoleApprover),
				Name:     "",
			},
			err: nil,
			preTest: func() {
				//monkey patch (replace) bcrypt.GenerateFromPassword with a function that returns a fixed output
				monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
					return []byte("testhash"), nil
				})

				// Generate "hash" to store from user password
				hash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
				domainUser := domain.User{Username: "testuser", Role: string(domain.RoleApprover), Password: string(hash)}
				userMock.EXPECT().GetByUsername("testuser").Return(domain.User{}, gorm.ErrRecordNotFound)
				createdUser := domainUser
				createdUser.ID = 1
				userMock.EXPECT().Save(domainUser).Return(createdUser, nil)
			},
		},
	}

	for _, test := range tests {
		test.preTest()
		if output, err := CreateUser(test.args); test.expected != output || err != test.err {
			t.Errorf("%s:\nArgs: %+v\nOutput: %+v\nExpected: %+v\n", test.id, test.args, output, test.expected)
		}
	}
}
