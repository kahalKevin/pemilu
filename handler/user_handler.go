package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"fmt"
	"bytes"
	"net/http"

	"service"
	"datasource"
	"repo"
	"restmodel"

	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/mux"
)

var db = datasource.InitConnection()
var userService = service.NewUserService(repo.NewRepository(db))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
		ID:        id,
		Name:      addUserRequest.Name,
		Tingkat:   addUserRequest.Tingkat,
		Username:  addUserRequest.Username,
		Password:  addUserRequest.Password,
		Role:      repo.CALON.String(),
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

func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	user, err := userService.ViewProfile(vars["usernameCalon"])

	if err != nil {
		log.Println("Cannot get User data", err)
		w.WriteHeader(http.StatusNotFound)
	}

	nameResponse := restmodel.ResponseGetUser {
		user.ID,
		user.Name,
		user.Tingkat,
	}
	json.NewEncoder(w).Encode(nameResponse)	
}

func AddPendukungHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	if "1" == r.Form["witness"][0]{
		isWitness = true
	}

	addPendukungRequest := restmodel.AddPendukungRequest {
		r.Form["idcalon"][0],
		r.Form["nik"][0],
		r.Form["firstname"][0],
		buf,
		r.Form["phone"][0],
		isWitness,
		handler.Filename,		
	}
	addResult, _ := userService.AddPendukung(addPendukungRequest, tokenHeader)
	var addResponse restmodel.ResponseGeneral
	addResponse.Result = addResult

	json.NewEncoder(w).Encode(addResponse)	
}
