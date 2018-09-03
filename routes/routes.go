package routes

import (
	"handler"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	route := mux.NewRouter()
	route.HandleFunc("/pemilu/login", handler.LoginHandler).Methods("POST")
	route.HandleFunc("/pemilu/addUser", handler.RegisterHandler).Methods("POST")
	// route.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	// route.HandleFunc("/viewprofile", handler.ProfileHandler).Methods("POST")
	return route
}
