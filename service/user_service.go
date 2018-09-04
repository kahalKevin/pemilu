package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"repo"
	"restmodel"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repo.UserRepository
}

type Token struct {
	jwt.StandardClaims
	ID       string `json:id`
	Username string `json:"username"`
	Role     string `json:"role"`
	Tingkat  string `json:"tingkat"`
}

type TokenData struct {
	Data  Token
	Token string
}

var mySigningKey []byte

func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

// NewUserService create new instance of UserService implementation
func NewUserService(userRepo repo.UserRepository) UserService {
	log.Println("NEW USER SERVICE")
	s := userService{userRepo: userRepo}
	return &s
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	log.Println(err)
	return err == nil
}

func (s *userService) Login(username string, password string) (tokenData TokenData, err error) {
	var token string
	mySigningKey := []byte("IDKWhatThisIs")

	userData, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
		} else {
			log.Println("Error at finding user's data", err)
		}
		return
	}

	match := CheckPasswordHash(password, userData.Password)
	if !match {
		err = errors.New("invalid password")
		log.Println("Wrong password")
		return
	}

	claims := Token{
		jwt.StandardClaims{
			Subject:   userData.ID,
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
		},
		userData.ID,
		userData.Username,
		userData.Role,
		userData.Tingkat,
	}

	signing := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ = signing.SignedString(mySigningKey)
	if len(token) == 0 {
		err = errors.New("Failed to generate token")
		log.Println("Failed to generate token")
		return
	}

	tokenData = TokenData{
		claims,
		token,
	}
	return
}

func (s *userService) Register(userRegister repo.User, token string) (registered bool, err error) {
	registered = false
	if len(token) == 0 || len(userRegister.Username) == 0 {
		err = errors.New("Must fill all field")
		log.Println("Fill all field", err)
		return
	}

	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}

	if dataToken.Role != repo.ADMIN.String() {
		err = errors.New("User dont have priviledges")
		return
	}

	checkUsername, err := s.userRepo.FindByUsername(userRegister.Username)
	if len(checkUsername.Username) != 0 {
		log.Println("Username exist on another account,    ", err)
		return
	}

	userRegister.Password, err = HashPassword(userRegister.Password)
	if err != nil {
		log.Println("Failed encrypting password,  ", err)
		return
	}

	_, err = s.userRepo.InsertNewUser(userRegister)
	if err != nil {
		log.Println("Failed registering,    ", err)
		return
	} else {
		registered = true
	}

	return
}

func validateToken(token string) (tokenData Token, err error) {
	at(time.Unix(0, 0), func() {
		tokenClaims, _ := jwt.ParseWithClaims(token, &Token{}, func(tokenClaims *jwt.Token) (interface{}, error) {
			return []byte("IDKWhatThisIs"), nil
		})

		if claims, _ := tokenClaims.Claims.(*Token); claims.ExpiresAt > time.Now().Unix() {
			id := claims.StandardClaims.Subject
			tokenData = Token{
				ID:       id,
				Username: claims.Username,
				Role:     claims.Role,
				Tingkat:  claims.Tingkat,
			}
		} else {
			err = errors.New("token Invalid")
		}
	})
	return
}

func (s *userService) ViewProfile(username string) (userProfile repo.User, err error) {
	userProfile, err = s.userRepo.FindByUsername(username)
	if err != nil {
		log.Println("Error at finding user's profile,	", err)
	}

	if "theboss" == userProfile.Username {
		userProfile = repo.User{}
		err = errors.New("Can not access this User")
	}
	return
}

func (s *userService) ConfirmDukungan(nik string, token string) (result bool, err error) {
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}
	tingkat := dataToken.Tingkat
	result, err = s.userRepo.ConfirmDukungan(nik, tingkat)
	if err != nil {
		log.Println("Error confirm dukungan,	", err)
	}
	return
}

func (s *userService) AddPendukung(request restmodel.AddPendukungRequest, token string) (success bool, err error) {
	var idCalon string
	var tingkat string
	var dataToken Token
	var errToken error
	autoConfirm := false
	success = false
	if "" == token {
		idCalon = request.IDCalon
		userProfile, _ := s.userRepo.FindByID(idCalon)
		if len(userProfile.ID) <= 0 {
			err = errors.New("Error at finding user's profile")
			log.Println("Error at finding user's profile,	", err)
			return
		}
		tingkat = userProfile.Tingkat
	} else {
		dataToken, errToken = validateToken(token)
		if errToken != nil {
			err = errToken
			return
		}
		idCalon = dataToken.ID
		tingkat = dataToken.Tingkat
		autoConfirm = true
	}

	dukungan, err := s.userRepo.FindAtDukungan(request.NIK, tingkat)
	if len(dukungan.ID) > 0 {
		err = errors.New("already Registered")
		return
	}

	var pendukung repo.Pendukung
	var newDukungan repo.Dukungan
	pendukung, err = s.userRepo.FindAtPendukung(request.NIK)
	if len(pendukung.ID) <= 0 {
		dataPendukung, _ := getSidalih3Data(request.NIK, request.Firstname)
		if len(dataPendukung.NIK) <= 0 {
			err = errors.New("NIK not registered at DPT")
			return
		}
		ID := uuid.Must(uuid.NewV4())
		extension, errImage := getImageExtension(request.FileName)
		if errImage != nil {
			err = errImage
			log.Println(err)
			return
		}
		generatedFileName := ID.String() + "." + extension
		go s.insertPendukung(dataPendukung, request, generatedFileName)
		go saveImage(request.Photo, generatedFileName)
	}

	nodeDukungan, errSFDukungan := snowflake.NewNode(1)
	if errSFDukungan != nil {
		err = errSFDukungan
		fmt.Println("Failed generating snowflake id,    ", err)
		return
	}
	idSFDukungan := nodeDukungan.Generate().String()

	newDukungan = repo.Dukungan{
		idSFDukungan,
		idCalon,
		request.NIK,
		tingkat,
		autoConfirm,
	}
	res, errInsert := s.userRepo.InsertDukungan(newDukungan)
	if err != nil {
		err = errInsert
		log.Println("Failed Insert Dukungan,    ", err)
		return
	}
	log.Println(newDukungan, res)
	success = true
	return
}

func saveImage(photo *bytes.Buffer, filename string) {
	file := photo.Bytes()
	f, err := os.OpenFile("./gbr/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	f.Write(file)
}

func (s *userService) insertPendukung(sidalih3Response restmodel.Sidalih3Response, request restmodel.AddPendukungRequest, fileName string) {
	gender := false
	if "L" == sidalih3Response.Gender {
		gender = true
	}
	nodePendukung, _ := snowflake.NewNode(2)
	idSFPendukung := nodePendukung.Generate().String()
	newPendukung := repo.Pendukung{
		idSFPendukung,
		sidalih3Response.Nama,
		sidalih3Response.NIK,
		sidalih3Response.Provinsi,
		sidalih3Response.Kabupaten,
		sidalih3Response.Kecamatan,
		sidalih3Response.Kelurahan,
		sidalih3Response.TPS,
		request.Phone,
		request.Witness,
		gender,
		fileName,
	}
	res, errInsert := s.userRepo.InsertPendukung(newPendukung)
	if errInsert != nil {
		log.Println("Failed Insert Pendukung,    ", errInsert)
	}
	log.Println(newPendukung, res)
}

func (s *userService) GetPendukungs(token string) (allPendukung restmodel.GetAllPendukungResponse, err error) {
	allPendukung.Data = make(map[string]restmodel.Site)
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}
	idCalon := dataToken.ID
	pendukungPart, errGetByCalon := s.userRepo.FindPendukungByCalon(idCalon)
	if errGetByCalon != nil {
		err = errGetByCalon
		return
	}
	for _, part := range pendukungPart {
		s := []string{part.Provinsi, part.Kabupaten, part.Kecamatan, part.Kelurahan, part.TPS}
		pendukung := restmodel.Pendukung{
			part.ID,
			part.Name,
			part.NIK,
			part.Phone,
			part.Witness,
			part.Gender,
			part.Status,
		}
		dataKey := strings.Join(s, ";")
		if dataValue, ok := allPendukung.Data[dataKey]; !ok {
			site := restmodel.Site{
				part.Provinsi,
				part.Kabupaten,
				part.Kecamatan,
				part.Kelurahan,
				part.TPS,
				[]restmodel.Pendukung{pendukung},
			}
			allPendukung.Data[dataKey] = site
		} else {
			listPendukung := dataValue.Pendukung
			listPendukung = append(listPendukung, pendukung)
			dataValue.Pendukung = listPendukung
			allPendukung.Data[dataKey] = dataValue
		}
	}
	return
}

func getSidalih3Data(nik string, name string) (sidalih3Response restmodel.Sidalih3Response, err error) {
	sidalih3Request := restmodel.Sidalih3Request{
		"search",
		nik,
		name,
	}
	reqJson, _ := json.Marshal(sidalih3Request)
	req, _ := http.NewRequest("POST", "https://sidalih3.kpu.go.id/dppublik/dpsnik", bytes.NewBuffer(reqJson))
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := &http.Client{}
	resp, errSidalih := client.Do(req)
	if errSidalih != nil {
		err = errSidalih
		log.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(body, &sidalih3Response)
	return
}

func getImageExtension(fileName string) (string, error) {
	stringSeparated := strings.Split(fileName, ".")
	lastElement := len(stringSeparated) - 1
	extension := make(map[string]bool)
	extension["jpg"] = true
	extension["png"] = true
	extension["jpeg"] = true

	if _, ok := extension[stringSeparated[lastElement]]; !ok {
		err := errors.New("extension Invalid")
		return "", err
	}

	return stringSeparated[lastElement], nil
}
