package repo

// UserRepository is a contract to persist user data to database
type UserRepository interface {
	// FindProfiles() ([]User, error)
	FindByID(id string) (User, error)
	FindAtDukungan(nik string, tingkat string) (Dukungan, error)
	FindAtPendukung(nik string) (Pendukung, error)
	InsertDukungan(dukungan Dukungan) (bool, error)
	InsertPendukung(pendukung Pendukung) (bool, error)
	// FindByEmail(email string) (User, error)
	// FindByMsisdn(msisdn string) (User, error)
	FindByUsername(usrname string) (User, error)
	FindPendukungByCalon(idCalon string) ([]PendukungPart, error)
	// FindUserRole(userID string) (UserRole, error)
	// FindExactRole(userID string, role int) (UserRole, error)
	InsertNewUser(user User) (string, error)
	// InsertToRole(userRole UserRole) (bool, error)
	ConfirmDukungan(nik string, tingkat string) (bool, error)
	DeleteDukungan(nik string, tingkat string) (bool, error)
}
