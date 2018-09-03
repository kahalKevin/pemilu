package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"service"

	"datasource"
	"repo"
	"restmodel"

	// "github.com/bwmarrin/snowflake"
)

var db = datasource.InitConnection()
var userService = service.NewUserService(repo.NewRepository(db))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 5000))

	var loginReq restmodel.LoginRequest
	err := json.Unmarshal(body, &loginReq)
	if err != nil {
		log.Println("ERROR at unmarshal", err)
		return
	}

	loginResult, err := userService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		log.Println("Failed at login,   ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var loginResp restmodel.Response

	if len(loginResult.Token) == 0 {
		loginResp.Result = false
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		loginResp.Result = true
		loginResp.Role = loginResult.Data.Role
		loginResp.Username = loginResult.Data.Username
		loginResp.Tingkat = loginResult.Data.Tingkat
	}

	w.Header().Set("token", loginResult.Token)
	json.NewEncoder(w).Encode(loginResp)
}

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 5000))

// 	var regRequest restmodel.RegisterRequest
// 	json.Unmarshal(body, &regRequest)

// 	node, err := snowflake.NewNode(1)
// 	if err != nil {
// 		log.Println("Fail to generate snowflake id,    ", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	id := node.Generate().String()

// 	userRegister := repo.User{
// 		ID:       id,
// 		Email:    regRequest.Email,
// 		Msisdn:   regRequest.Msisdn,
// 		Username: regRequest.Username,
// 		Password: regRequest.Password,
// 		Status:   0,
// 	}

// 	role := regRequest.Role

// 	registerResult, err := userService.Register(userRegister, role)
// 	if err != nil {
// 		log.Println("failed to register,    ", err)
// 		w.WriteHeader(http.StatusNotAcceptable)
// 	}

// 	var regResponse restmodel.Response

// 	if !registerResult {
// 		regResponse.Message = "Register failed"
// 	} else {
// 		regResponse.Message = "Register success"
// 		json.NewEncoder(w).Encode(regResponse)
// 		w.WriteHeader(http.StatusBadRequest)
// 	}

// 	json.NewEncoder(w).Encode(regResponse)
// }

// func ProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	tokenHeader := r.Header.Get("token")

// 	profile, err := userService.ViewProfile(tokenHeader)
// 	if err != nil {
// 		log.Println("Failed to view profile,    ", err)
// 		w.WriteHeader(401)
// 	}

// 	profile.Password = "secret"

// 	json.NewEncoder(w).Encode(profile)
// }
