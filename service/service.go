package service

import (
	"repo"
)

// UserService will be implemented in user_service
type UserService interface {
	Login(username string, password string) (TokenData, error)
	Register(userRegister repo.User, token string) (bool, error)
	// ViewProfile(token string) (repo.User, error)
}

var User UserService

func NewService(userRepo repo.UserRepository) {
	User = NewUserService(userRepo)
}
