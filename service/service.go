package service

import (
	"repo"
	"restmodel"
)

// UserService will be implemented in user_service
type UserService interface {
	Login(username string, password string) (TokenData, error)
	Register(userRegister restmodel.AddUserRequest, token string) (bool, error)
	ViewProfile(username string) (repo.User, error)
	AddPendukung(request restmodel.AddPendukungRequest, token string) (bool, error)
	insertPendukung(sidalih3Response restmodel.Sidalih3Response, request restmodel.AddPendukungRequest, fileName string)
	GetPendukungs(token string) (restmodel.GetAllPendukungResponse, error)
	ConfirmDukungan(nik string, token string) (bool, error)
	DeleteDukungan(nik string, token string) (bool, error)
	GetPendukungFull(nik string, token string) (repo.PendukungFull, error)
	GetUsers(token string) ([]repo.UserPart, error)
	ChangePassword(token string, password string, newPassword string) (bool, error)
	DeleteUser(idCalon string, token string) (bool, error)
}

var User UserService

func NewService(userRepo repo.UserRepository) {
	User = NewUserService(userRepo)
}
