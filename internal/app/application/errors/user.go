package errors

import "net/http"

type UserNotFoundError struct{}

func (e UserNotFoundError) Error() string {
	return "user not found"
}

func (e UserNotFoundError) Code() int {
	return http.StatusNotFound
}

type UserRoleInvalidError struct{}

func (e UserRoleInvalidError) Error() string {
	return "user role invalid"
}

func (e UserRoleInvalidError) Code() int {
	return http.StatusNotAcceptable
}

type UsernameAlreadyExistsError struct{}

func (e UsernameAlreadyExistsError) Error() string {
	return "username already exists"
}

func (e UsernameAlreadyExistsError) Code() int {
	return http.StatusConflict
}

type InvalidPasswordStringError struct{}

func (e InvalidPasswordStringError) Error() string {
	return "password string is invalid"
}

func (e InvalidPasswordStringError) Code() int {
	return http.StatusBadRequest
}

type UserCreationFailed struct{}

func (e UserCreationFailed) Error() string {
	return "user creation failed"
}

func (e UserCreationFailed) Code() int {
	return http.StatusInternalServerError
}
