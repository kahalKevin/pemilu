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
	insertPendukung(sidalih3Response restmodel.Sidalih3Response, request restmodel.AddPendukungRequest, fileName string)
	GetPendukungs(token string) (restmodel.GetAllPendukungResponse, error)
	ConfirmDukungan(nik string, token string) (bool, error)
	DeleteDukungan(nik string, token string) (bool, error)
	GetPendukungFull(nik string, token string) (repo.PendukungFull, error)
}

var User UserService

func NewService(userRepo repo.UserRepository) {
	User = NewUserService(userRepo)
}
