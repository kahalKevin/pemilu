package restmodel

import "bytes"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddUserRequest struct {
	Name     string `json:"name"`
	Tingkat  string `json:"tingkat"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddPendukungRequest struct {
	IDCalon   string        `json:"idcalon"`
	NIK       string        `json:"nik"`
	Firstname string        `json:"firstname"`
	Photo     *bytes.Buffer `json:"photo"`
	Phone     string        `json:"phone"`
	Witness   bool          `json:"witness"`
	FileName  string
}

type Response struct {
	Result   bool   `json:"result"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Tingkat  string `json:"tingkat"`
}

type ResponseGetUser struct {
	IDCalon string `json:"idCalon"`
	Name    string `json:"name"`
	Tingkat string `json:"tingkat"`
}

type ResponseGeneral struct {
	Result bool `json:"result"`
}

type RegisterRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Msisdn   string `json:"msisdn"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Role     int    `json:"role"`
}

type Sidalih3Request struct {
	Command string `json:"cmd"`
	NIK     string `json:"nik"`
	Nama    string `json:"nama"`
}

type Sidalih3Response struct {
	Nama      string `json:"nama"`
	NIK       string `json:"nik"`
	TPS       string `json:"tps"`
	Gender    string `json:"jenis_kelamin"`
	Kelurahan string `json:"kelurahan"`
	Kecamatan string `json:"kecamatan"`
	Kabupaten string `json:"kabupaten"`
	Provinsi  string `json:"provinsi"`
}

type Pendukung struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	NIK     string `json:"nik"`
	Phone   string `json:"phone"`
	Witness bool   `json:"witness"`
	Gender  bool   `json:"gender"`
	Status  bool   `json:"status"`
}

type Site struct {
	Provinsi  string      `json:"provinsi"`
	Kabupaten string      `json:"kabupaten"`
	Kecamatan string      `json:"kecamatan"`
	Kelurahan string      `json:"kelurahan"`
	TPS       string      `json:"tps"`
	Pendukung []Pendukung `json:"pendukungs"`
}

type GetAllPendukungResponse struct {
	Data map[string]Site `json:"data"`
}
