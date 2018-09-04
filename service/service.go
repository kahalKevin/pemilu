package service

import (
	"repo"
	"restmodel"
)

// UserService will be implemented in user_service
type UserService interface {
	Login(username string, password string) (TokenData, error)
	Register(userRegister repo.User, token string) (bool, error)
	ViewProfile(username string) (repo.User, error)
	AddPendukung(request restmodel.AddPendukungRequest, token string) (bool, error)
}

var User UserService

func NewService(userRepo repo.UserRepository) {
	User = NewUserService(userRepo)
}
