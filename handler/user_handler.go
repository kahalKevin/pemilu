package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"datasource"
	"repo"
	"restmodel"
	"service"

	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/mux"
)

var db = datasource.InitConnection()
var userService = service.NewUserService(repo.NewRepository(db))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 5000))

	var loginResp restmodel.Response
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
		loginResp.Result = false
		json.NewEncoder(w).Encode(loginResp)
		return
	}

	if len(loginResult.Token) == 0 {
		loginResp.Result = false
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		loginResp.Result = true
		loginResp.Role = loginResult.Data.Role
		loginResp.Username = loginResult.Data.Username
		loginResp.Tingkat = loginResult.Data.Tingkat
		loginResp.AvatarUrl = loginResult.AvatarUrl
	}

	w.Header().Set("Access-Control-Expose-Headers", "token")
	w.Header().Set("token", loginResult.Token)
	json.NewEncoder(w).Encode(loginResp)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusOK)

	tokenHeader := r.Header.Get("token")
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 5000))

	var addUserRequest restmodel.AddUserRequest
	json.Unmarshal(body, &addUserRequest)

	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Println("Fail to generate snowflake id,    ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := node.Generate().String()

	userRegister := repo.User{
		ID:       id,
		Name:     addUserRequest.Name,
		Tingkat:  addUserRequest.Tingkat,
		Username: addUserRequest.Username,
		Password: addUserRequest.Password,
		Role:     repo.CALON.String(),
	}

	registerResult, err := userService.Register(userRegister, tokenHeader)
	if err != nil {
		log.Println("failed to register,    ", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var regResponse restmodel.ResponseGeneral
	regResponse.Result = registerResult

	json.NewEncoder(w).Encode(regResponse)
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 5000))
	var changePasswordReq restmodel.ChangePasswordRequest
	err := json.Unmarshal(body, &changePasswordReq)
	if err != nil {
		log.Println("ERROR at unmarshal", err)
		return
	}
	result, _ := userService.ChangePassword(tokenHeader, changePasswordReq.OldPassword, changePasswordReq.NewPassword)
	var addResponse restmodel.ResponseGeneral
	addResponse.Result = result
	json.NewEncoder(w).Encode(addResponse)
}

func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)

	user, err := userService.ViewProfile(vars["usernameCalon"])

	if err != nil {
		log.Println("Cannot get User data", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	nameResponse := restmodel.ResponseGetUser{
		user.ID,
		user.Name,
		user.Tingkat,
		user.AvatarUrl,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nameResponse)
}

func AddPendukungHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	var threshold int64
	// 1 << 19 to make 512kb
	threshold = 1 << 19
	r.ParseMultipartForm(threshold)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	if handler.Size > threshold {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	isWitness := false
	if "1" == r.Form["witness"][0] {
		isWitness = true
	}

	addPendukungRequest := restmodel.AddPendukungRequest{
		r.Form["idcalon"][0],
		r.Form["nik"][0],
		r.Form["firstname"][0],
		buf,
		r.Form["phone"][0],
		isWitness,
		r.Form["address"][0],
		handler.Filename,
	}
	addResult, _ := userService.AddPendukung(addPendukungRequest, tokenHeader)
	var addResponse restmodel.ResponseGeneral
	addResponse.Result = addResult

	json.NewEncoder(w).Encode(addResponse)
}

func GetPendukungsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	user, _ := userService.GetPendukungs(tokenHeader)
	json.NewEncoder(w).Encode(user)
}

func ConfirmPendukungHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	niks, _ := r.URL.Query()["nik"]
	nik := niks[0]
	result, _ := userService.ConfirmDukungan(nik, tokenHeader)
	var addResponse restmodel.ResponseGeneral
	addResponse.Result = result
	json.NewEncoder(w).Encode(addResponse)
}

func GetPendukungHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	niks, _ := r.URL.Query()["nik"]
	nik := niks[0]
	fullData, _ := userService.GetPendukungFull(nik, tokenHeader)
	json.NewEncoder(w).Encode(fullData)
}

func DeletePendukungHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	niks, ok := r.URL.Query()["nik"]
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	nik := niks[0]
	result, _ := userService.DeleteDukungan(nik, tokenHeader)
	var delResponse restmodel.ResponseGeneral
	delResponse.Result = result
	json.NewEncoder(w).Encode(delResponse)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	fullData, _ := userService.GetUsers(tokenHeader)
	json.NewEncoder(w).Encode(fullData)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}

	tokenHeader := r.Header.Get("token")
	idCalon, ok := r.URL.Query()["id"]
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	id := idCalon[0]
	result, _ := userService.DeleteUser(id, tokenHeader)
	var delResponse restmodel.ResponseGeneral
	delResponse.Result = result
	json.NewEncoder(w).Encode(delResponse)
}

func setupResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, token")
}
