package routes

import (
	"handler"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	route := mux.NewRouter()
	route.HandleFunc("/pemilu/login", handler.LoginHandler).Methods("POST")
	route.HandleFunc("/pemilu/addUser", handler.RegisterHandler).Methods("POST")
	route.HandleFunc("/pemilu/addPendukung", handler.AddPendukungHandler).Methods("POST")

	route.HandleFunc("/pemilu/changePassword", handler.ChangePasswordHandler).Methods("POST")

	route.HandleFunc("/pemilu/getUsers", handler.GetUsersHandler).Methods("GET")
	route.HandleFunc("/pemilu/getPendukungs", handler.GetPendukungsHandler).Methods("GET")
	route.HandleFunc("/pemilu/confirmPendukung", handler.ConfirmPendukungHandler).Methods("GET")
	route.HandleFunc("/pemilu/getPendukung", handler.GetPendukungHandler).Methods("GET")
	route.HandleFunc("/pemilu/deletePendukung", handler.DeletePendukungHandler).Methods("GET")
	route.HandleFunc("/pemilu/{usernameCalon}", handler.GetNameHandler).Methods("GET")
	return route
}
