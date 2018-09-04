package repo

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	conn *sqlx.DB
	// findAllStmt       *sqlx.Stmt
	findIDStmt      *sqlx.Stmt
	findAtDukungan  *sqlx.Stmt
	findAtPendukung *sqlx.Stmt
	insertDukungan  *sqlx.NamedStmt
	insertPendukung *sqlx.NamedStmt
	// findEmailStmt     *sqlx.Stmt
	// findMsisdnStmt    *sqlx.Stmt
	findUsrnameStmt      *sqlx.Stmt
	findPendukungByCalon *sqlx.Stmt
	// findRoleStmt      *sqlx.Stmt
	// findExactRoleStmt *sqlx.Stmt
	insertUser *sqlx.NamedStmt
	// insertToRole      *sqlx.NamedStmt
}

func (db *userRepository) MustPrepareStmt(query string) *sqlx.Stmt {
	stmt, err := db.conn.Preparex(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %s\n", err)
		os.Exit(1)
	}
	return stmt
}

func (db *userRepository) MustPrepareNamedStmt(query string) *sqlx.NamedStmt {
	stmt, err := db.conn.PrepareNamed(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %s\n", err)
		os.Exit(1)
	}
	return stmt
}

// NewRepository create new instance of UserRepository implementation
func NewRepository(db *sqlx.DB) UserRepository {
	log.Println("NEW USER REPOSITORY")
	r := userRepository{conn: db}
	// r.findAllStmt = r.MustPrepareStmt("SELECT * FROM user_auth")
	r.findIDStmt = r.MustPrepareStmt("SELECT * FROM User WHERE id=?")
	r.findAtDukungan = r.MustPrepareStmt("SELECT * FROM Dukungan WHERE nik =? AND tingkat=?")
	r.findAtPendukung = r.MustPrepareStmt("SELECT * FROM Pendukung WHERE nik=?")
	// r.findMsisdnStmt = r.MustPrepareStmt("SELECT * FROM user_auth WHERE msisdn=?")
	// r.findEmailStmt = r.MustPrepareStmt("SELECT * FROM user_auth WHERE email=?")
	r.findUsrnameStmt = r.MustPrepareStmt("SELECT * FROM User WHERE username=?")
	r.findPendukungByCalon = r.MustPrepareStmt("select p.id, p.nik, name, phone, witness, gender, status, provinsi, kabupaten, kecamatan, kelurahan, tps from Pendukung p inner join Dukungan d on p.nik=d.nik where idCalon=?")
	// r.findRoleStmt = r.MustPrepareStmt("SELECT * FROM user_role WHERE user_id =?")
	// r.findExactRoleStmt = r.MustPrepareStmt("SELECT * FROM user_role WHERE user_id =? AND role=?")
	r.insertUser = r.MustPrepareNamedStmt("INSERT INTO User (id, name, tingkat, username, password, role) VALUES (:id, :name, :tingkat, :username, :password, :role)")
	r.insertDukungan = r.MustPrepareNamedStmt("INSERT INTO Dukungan (id, idCalon, nik, tingkat, status) VALUES (:id, :idCalon, :nik, :tingkat, :status)")
	r.insertPendukung = r.MustPrepareNamedStmt("INSERT INTO Pendukung (id, name, nik, provinsi, kabupaten, kecamatan, kelurahan, tps, phone, witness, gender, photo) VALUES (:id, :name, :nik, :provinsi, :kabupaten, :kecamatan, :kelurahan, :tps, :phone, :witness, :gender, :photo)")
	// r.insertToRole = r.MustPrepareNamedStmt("INSERT INTO user_role (id, user_id, role) VALUES (:id, :user_id, :role)")
	return &r
}

// func (db *userRepository) FindProfiles() (usr []User, err error) {
// 	err = db.findAllStmt.Select(&usr)
// 	if err != nil {
// 		log.Println("Error at finding profiles,    ", err)
// 	}
// 	return
// }

func (db *userRepository) FindByID(id string) (usr User, err error) {
	err = db.findIDStmt.Get(&usr, id)
	if err != nil {
		log.Printf("ID: %v , doesn't exist", id)
		log.Println(err)
	}
	return
}

// func (db *userRepository) FindByMsisdn(msisdn string) (usr User, err error) {
// 	err = db.findMsisdnStmt.Get(&usr, msisdn)
// 	if err != nil {
// 		log.Println("Error at Finding phone number,    ", err)
// 	}
// 	return
// }

// func (db *userRepository) FindByEmail(email string) (usr User, err error) {
// 	var user []User
// 	err = db.findEmailStmt.Select(&user, email)
// 	if err != nil {
// 		log.Println("Error at finding email,    ", err)
// 	}

// 	if len(user) != 0 {
// 		usr = user[0]
// 	}

// 	return
// }

func (db *userRepository) FindByUsername(usrname string) (usr User, err error) {
	var u User
	err = db.findUsrnameStmt.Get(&u, usrname)
	usr = u
	if err != nil {
		log.Println("Error at finding username,    ", err)
	}
	return
}

func (db *userRepository) FindPendukungByCalon(idCalon string) (pendukungPart []PendukungPart, err error) {
	err = db.findPendukungByCalon.Select(&pendukungPart, idCalon)
	if err != nil {
		log.Println("Error at finding pendukung by calon,    ", err)
	}
	return
}

// func (db *userRepository) FindUserRole(userID string) (userRole UserRole, err error) {
// 	err = db.findRoleStmt.Get(&userRole, userID)
// 	if err != nil {
// 		log.Println("Error while finding user role,    ", err)
// 	}
// 	return
// }

func (db *userRepository) FindAtPendukung(nik string) (pendukung Pendukung, err error) {
	err = db.findAtPendukung.Get(&pendukung, nik)
	if err != nil {
		log.Println("Error at finding row at Pendukung,    ", err)
	}
	return
}

func (db *userRepository) FindAtDukungan(nik string, tingkat string) (dukungan Dukungan, err error) {
	err = db.findAtDukungan.Get(&dukungan, nik, tingkat)
	if err != nil {
		log.Println("Error at finding row at Dukungan,    ", err)
	}
	return
}

func (db *userRepository) InsertDukungan(dukungan Dukungan) (success bool, err error) {
	_, err = db.insertDukungan.Exec(dukungan)
	if err != nil {
		log.Println("Error inserting new Dukungan,    ", err)
	} else {
		success = true
	}
	return
}

func (db *userRepository) InsertPendukung(pendukung Pendukung) (success bool, err error) {
	_, err = db.insertPendukung.Exec(pendukung)
	if err != nil {
		log.Println("Error inserting new Pendukung,    ", err)
	} else {
		success = true
	}
	return
}

func (db *userRepository) InsertNewUser(user User) (lastID string, err error) {
	_, err = db.insertUser.Exec(user)
	if err != nil {
		log.Println("Error inserting new user,    ", err)
	}

	lastID = user.ID

	return
}

// func (db *userRepository) InsertToRole(newRole UserRole) (success bool, err error) {
// 	_, err = db.insertToRole.Exec(newRole)
// 	if err != nil {
// 		log.Println("Error inserting new role,    ", err)
// 	} else {
// 		success = true
// 	}
// 	return
// }
