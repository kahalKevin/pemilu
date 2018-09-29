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
	route.HandleFunc("/pemilu/deleteUser", handler.DeleteUserHandler).Methods("DELETE")
	route.HandleFunc("/pemilu/getUsers", handler.GetUsersHandler).Methods("GET")
	route.HandleFunc("/pemilu/getPendukungs", handler.GetPendukungsHandler).Methods("GET")
	route.HandleFunc("/pemilu/confirmPendukung", handler.ConfirmPendukungHandler).Methods("GET")
	route.HandleFunc("/pemilu/getPendukung", handler.GetPendukungHandler).Methods("GET")
	route.HandleFunc("/pemilu/deletePendukung", handler.DeletePendukungHandler).Methods("DELETE")
	route.HandleFunc("/pemilu/{usernameCalon}", handler.GetNameHandler).Methods("GET")

	route.HandleFunc("/pemilu/login", handler.LoginHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/addUser", handler.RegisterHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/addPendukung", handler.AddPendukungHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/changePassword", handler.ChangePasswordHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/deleteUser", handler.DeleteUserHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/getUsers", handler.GetUsersHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/getPendukungs", handler.GetPendukungsHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/confirmPendukung", handler.ConfirmPendukungHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/getPendukung", handler.GetPendukungHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/deletePendukung", handler.DeletePendukungHandler).Methods("OPTIONS")
	route.HandleFunc("/pemilu/{usernameCalon}", handler.GetNameHandler).Methods("OPTIONS")
	return route
}
