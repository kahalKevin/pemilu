package service

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"repo"
	"restmodel"
	"strconv"
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
	Data      Token
	Token     string
	AvatarUrl string
}

var mySigningKey []byte

func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

const urlImg string = "http://solagratia.web.id/images/"

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
		urlImg + userData.AvatarUrl,
	}
	return
}

func (s *userService) Register(userRegister restmodel.AddUserRequest, token string) (registered bool, err error) {
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

	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Println("Fail to generate snowflake id,    ", err)
		return
	}

	ID := uuid.Must(uuid.NewV4())
	extension, errImage := getImageExtension(userRegister.FileName)
	if errImage != nil {
		err = errImage
		log.Println(err)
		return
	}
	generatedFileName := ID.String() + "." + extension

	id := node.Generate().String()
	registerData := repo.User{
		ID:        id,
		Name:      userRegister.Name,
		Tingkat:   userRegister.Tingkat,
		Username:  userRegister.Username,
		Password:  userRegister.Password,
		Role:      repo.CALON.String(),
		AvatarUrl: generatedFileName,
	}

	_, err = s.userRepo.InsertNewUser(registerData)

	if err != nil {
		log.Println("Failed registering,    ", err)
		return
	} else {
		registered = true
		go saveImage(userRegister.AvatarUrl, generatedFileName)
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
	} else {
		userProfile.AvatarUrl = urlImg + userProfile.AvatarUrl
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

func (s *userService) ChangePassword(token string, password string, newPassword string) (result bool, err error) {
	result = false
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}

	username := dataToken.Username

	user, _ := s.userRepo.FindByUsername(username)
	if len(user.Username) == 0 {
		log.Println("Username not exist")
		return
	}

	valid := CheckPasswordHash(password, user.Password)
	if valid {
		newPassword, _ = HashPassword(newPassword)
		result, err = s.userRepo.ChangePassword(username, newPassword)
		if err != nil {
			log.Println("Error change password,	", err)
		}
	}
	return
}

func (s *userService) DeleteDukungan(nik string, token string) (result bool, err error) {
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}
	tingkat := dataToken.Tingkat
	result, err = s.userRepo.DeleteDukungan(nik, tingkat)
	if err != nil {
		log.Println("Error delete dukungan,	", err)
	}
	return
}

func (s *userService) GetPendukungFull(nik string, token string) (full repo.PendukungFull, err error) {
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}
	tingkat := dataToken.Tingkat
	full, err = s.userRepo.GetPendukungFull(nik, tingkat)
	if err != nil {
		log.Println("Error get pendukung full,	", err)
	}
	full.Photo = urlImg + full.Photo
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
		dataPendukung, _ := pemiluDataAggregate(request.NIK, request.Firstname)
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
		// go saveImage(request.Photo, generatedFileName)
	}

	nodeDukungan, errSFDukungan := snowflake.NewNode(1)
	if errSFDukungan != nil {
		err = errSFDukungan
		fmt.Println("Failed generating snowflake id,    ", err)
		return
	}
	idSFDukungan := nodeDukungan.Generate().String()
	t := time.Now().Unix()
	timestamp := strconv.FormatInt(t, 10)

	newDukungan = repo.Dukungan{
		idSFDukungan,
		idCalon,
		request.NIK,
		tingkat,
		autoConfirm,
		timestamp,
	}
	_, errInsert := s.userRepo.InsertDukungan(newDukungan)
	if errInsert != nil {
		err = errInsert
		log.Println("Failed Insert Dukungan,    ", err)
		return
	}
	success = true
	return
}

func saveImage(photo *bytes.Buffer, filename string) {
	file := photo.Bytes()
	f, err := os.OpenFile("/var/www/web-sola-gratia-yii2/backend/web/images/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
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
		request.Address,
	}
	_, errInsert := s.userRepo.InsertPendukung(newPendukung)
	if errInsert != nil {
		log.Println("Failed Insert Pendukung,    ", errInsert)
	}
}

func (s *userService) GetUsers(token string) (users []repo.UserPart, err error) {
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}

	if dataToken.Role != repo.ADMIN.String() {
		err = errors.New("User dont have priviledges")
		return
	}

	users, err = s.userRepo.GetUsers()
	return
}

func (s *userService) DeleteUser(idCalon string, token string) (result bool, err error) {
	result = false
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}
	if dataToken.Role != repo.ADMIN.String() {
		err = errors.New("User dont have priviledges")
		return
	}
	s.userRepo.DeleteDukunganByCalon(idCalon)
	s.userRepo.DeleteUser(idCalon)
	result = true
	return
}

func (s *userService) GetPendukungs(token, start, end string) (allPendukung restmodel.GetAllPendukungResponse, err error) {
	allPendukung.Data = make(map[string]restmodel.Site)
	dataToken, errToken := validateToken(token)
	if errToken != nil {
		err = errToken
		return
	}
	idCalon := dataToken.ID

	var pendukungPart []repo.PendukungPart
	var errGetByCalon error
	if len(start) < 1 {
		pendukungPart, errGetByCalon = s.userRepo.FindPendukungByCalon(idCalon)
	} else {
		pendukungPart, errGetByCalon = s.userRepo.FindPendukungByCalonAndLimit(idCalon, start, end)
	}

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
			part.Timestamp,
			part.Address,
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

	clientReal := &http.Client{}
	client := &http.Client{
		CheckRedirect: func() func(req *http.Request, via []*http.Request) error {
			redirects := 0
			return func(req *http.Request, via []*http.Request) error {
				if redirects > 1 {
					return errors.New("stopped after 1 redirects")
				}
				redirects++
				return nil
			}
		}(),
	}

	resp, errSidalih := client.Do(req)
	if errSidalih != nil {
		var cookieName string
		for _, cookie := range resp.Cookies() {
			cookieName = cookie.Name
			cookie := http.Cookie{Name: cookie.Name, Value: cookie.Value}
			req.AddCookie(&cookie)
		}
		egovCookie, _ := req.Cookie(cookieName)
		if egovCookie == nil {
			err = errSidalih
			log.Println(err)
			return
		}
	}
	resp, errSidalih = clientReal.Do(req)
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
	if "" == fileName {
		err := errors.New("extension Invalid")
		return "", err
	}
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

func pemiluDataAggregate(nik string, name string) (dataResponse restmodel.Sidalih3Response, err error) {
	fail := make(chan uint32, 2)
	success := make(chan restmodel.Sidalih3Response, 1)

	go getLPHMquick(nik, fail, success)
	go getLHPMdata(nik, name, fail, success)
	go getSidalih3DataV2(nik, name, fail, success)

	var failCount uint32
	failCount = 0
	for {
		select {
		case doom := <-fail:
			failCount = failCount + doom
			if failCount >= 3 {
				err = errors.New("fail call LHPM")
				return
			}
		case bless := <-success:
			dataResponse = bless
			return
		}
	}
}

func getLPHMquick(nik string, fail chan<- uint32, success chan<- restmodel.Sidalih3Response) {
	defer elapsed("getLPHMquick")()
	defer recoverKpuCall("kmbmicro", fail)
	var dataResponse restmodel.Sidalih3Response
	var err error
	apiUrl := "https://kmbmicro.xyz"
	resource := "experiment/pemilu2019.php/"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}

	req, _ := http.NewRequest("GET", urlStr, nil)

	q := req.URL.Query()
	q.Add("nik", nik)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fail <- 1
		return
	}
	defer resp.Body.Close()
	var respData restmodel.LindungiHPMResponse
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &respData)
	if respData.Message != "success" {
		fail <- 1
		return
	} else {
		dataResponse.NIK = nik
		dataResponse.Nama = respData.Data.Nama
		dataResponse.TPS = respData.Data.TPS
		dataResponse.Gender = respData.Data.Sex
		dataResponse.Kelurahan = respData.Data.Kelurahan
		dataResponse.Kecamatan = respData.Data.Kecamatan
		dataResponse.Kabupaten = respData.Data.KabKota
		dataResponse.Provinsi = respData.Data.Provinsi
		fmt.Println("GET FROM kmbmicro.xyz")
	}
	success <- dataResponse
}

func getLHPMdata(nik string, name string, fail chan<- uint32, success chan<- restmodel.Sidalih3Response) {
	defer elapsed("getLHPMdata")()
	defer recoverKpuCall("lindungihakpilihmu", fail)
	var dataResponse restmodel.Sidalih3Response
	var err error
	apiUrl := "https://lindungihakpilihmu.kpu.go.id"
	resource := "/index.php/dpt/proses_ceknik/"
	data := url.Values{}
	data.Set("nik", nik)
	data.Add("nama", name)

	failedSign := "message:failed"
	dataSign := "data:{"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fail <- 1
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	bodyString = strings.Replace(bodyString, "\\\"", "", -1)
	bodyString = strings.Replace(bodyString, "\\", "", -1)
	bodyString = strings.Replace(bodyString, "\"", "", -1)
	if strings.Contains(bodyString, failedSign) {
		fail <- 1
		return
	} else {
		i := strings.Index(bodyString, dataSign) + 6
		data := bodyString[i : len(bodyString)-2]
		dataArray := strings.Split(data, ",")
		dataResponse.NIK = nik
		for _, part := range dataArray {
			content := strings.Split(part, ":")
			if "nama" == content[0] {
				dataResponse.Nama = content[1]
			} else if "tps" == content[0] {
				dataResponse.TPS = content[1]
			} else if "jenis_kelamin" == content[0] {
				dataResponse.Gender = content[1]
			} else if "namaKelurahan" == content[0] {
				dataResponse.Kelurahan = content[1]
			} else if "namaKecamatan" == content[0] {
				dataResponse.Kecamatan = content[1]
			} else if "namaKabKota" == content[0] {
				dataResponse.Kabupaten = content[1]
			} else if "namaPropinsi" == content[0] {
				dataResponse.Provinsi = content[1]
			}
		}
		fmt.Println("GET FROM lindungihakpilihmu.kpu.go.id")
	}
	success <- dataResponse
}

func getSidalih3DataV2(nik string, name string, fail chan<- uint32, success chan<- restmodel.Sidalih3Response) {
	defer elapsed("getSidalih3DataV2")()
	defer recoverKpuCall("sidalih3", fail)
	sidalih3Request := restmodel.Sidalih3Request{
		"search",
		nik,
		name,
	}
	reqJson, _ := json.Marshal(sidalih3Request)
	req, _ := http.NewRequest("POST", "https://sidalih3.kpu.go.id/dppublik/dpsnik", bytes.NewBuffer(reqJson))
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, errSidalih := client.Do(req)
	if errSidalih != nil {
		fail <- 1
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var sidalih3Response restmodel.Sidalih3Response
	json.Unmarshal(body, &sidalih3Response)
	if len(sidalih3Response.NIK) <= 0 {
		fail <- 1
		return
	}
	success <- sidalih3Response
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

func recoverKpuCall(funcName string, fail chan<- uint32) {
	r := recover()
	if r != nil {
		fmt.Println("Panic on Call" + funcName)
		fail <- 1
	}
}
