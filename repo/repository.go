package repo

// UserRepository is a contract to persist user data to database
type UserRepository interface {
	// FindProfiles() ([]User, error)
	FindByID(id string) (User, error)
	FindAtDukungan(nik string, tingkat string) (Dukungan, error)
	FindAtPendukung(nik string) (Pendukung, error)
	InsertDukungan(dukungan Dukungan) (bool, error)
	InsertPendukung(pendukung Pendukung) (bool, error)
	FindByUsername(usrname string) (User, error)
	FindPendukungByCalon(idCalon string) ([]PendukungPart, error)
	InsertNewUser(user User) (string, error)
	ConfirmDukungan(nik string, tingkat string) (bool, error)
	DeleteDukungan(nik string, tingkat string) (bool, error)
	GetPendukungFull(nik string, tingkat string) (PendukungFull, error)
	GetUsers() ([]UserPart, error)
}
