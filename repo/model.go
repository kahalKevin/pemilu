package repo

type User struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Tingkat   string `json:"tingkat" db:"tingkat"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Role      string `json:"role" db:"role"`
	AvatarUrl string `json:"avatarUrl" db:"avatarUrl"`
}

type Pendukung struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	NIK       string `json:"nik" db:"nik"`
	Provinsi  string `json:"provinsi" db:"provinsi"`
	Kabupaten string `json:"kabupaten" db:"kabupaten"`
	Kecamatan string `json:"kecamatan" db:"kecamatan"`
	Kelurahan string `json:"kelurahan" db:"kelurahan"`
	TPS       string `json:"tps" db:"tps"`
	Phone     string `json:"phone" db:"phone"`
	Witness   bool   `json:"witness" db:"witness"`
	Gender    bool   `json:"gender" db:"gender"`
	Photo     string `json:"photo" db:"photo"`
	Address   string `json:"address" db:"address"`
}

type PendukungFull struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	NIK       string `json:"nik" db:"nik"`
	Provinsi  string `json:"provinsi" db:"provinsi"`
	Kabupaten string `json:"kabupaten" db:"kabupaten"`
	Kecamatan string `json:"kecamatan" db:"kecamatan"`
	Kelurahan string `json:"kelurahan" db:"kelurahan"`
	TPS       string `json:"tps" db:"tps"`
	Phone     string `json:"phone" db:"phone"`
	Witness   bool   `json:"witness" db:"witness"`
	Gender    bool   `json:"gender" db:"gender"`
	Photo     string `json:"photo" db:"photo"`
	Status    bool   `json:"status" db:"status"`
	Address   string `json:"address" db:"address"`
}

type PendukungPart struct {
	ID        string `json:"id" db:"id"`
	NIK       string `json:"nik" db:"nik"`
	Name      string `json:"name" db:"name"`
	Phone     string `json:"phone" db:"phone"`
	Witness   bool   `json:"witness" db:"witness"`
	Gender    bool   `json:"gender" db:"gender"`
	Status    bool   `json:"status" db:"status"`
	Provinsi  string `json:"provinsi" db:"provinsi"`
	Kabupaten string `json:"kabupaten" db:"kabupaten"`
	Kecamatan string `json:"kecamatan" db:"kecamatan"`
	Kelurahan string `json:"kelurahan" db:"kelurahan"`
	TPS       string `json:"tps" db:"tps"`
	Timestamp string `json:"timestamp" db:"timestamp"`
	Address   string `json:"address" db:"address"`
}

type Dukungan struct {
	ID        string `json:"id" db:"id"`
	IDCalon   string `json:"idCalon" db:"idCalon"`
	NIK       string `json:"nik" db:"nik"`
	Tingkat   string `json:"tingkat" db:"tingkat"`
	Status    bool   `json:"status" db:"status"`
	Timestamp string `json:"timestamp" db:"timestamp"`
}

type UserPart struct {
	ID      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Tingkat string `json:"tingkat" db:"tingkat"`
}
