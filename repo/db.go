package repo

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	conn                  *sqlx.DB
	findIDStmt            *sqlx.Stmt
	findAtDukungan        *sqlx.Stmt
	findAtPendukung       *sqlx.Stmt
	insertDukungan        *sqlx.NamedStmt
	insertPendukung       *sqlx.NamedStmt
	findUsrnameStmt       *sqlx.Stmt
	findPendukungByCalon  *sqlx.Stmt
	insertUser            *sqlx.NamedStmt
	confirmDukungan       *sqlx.Stmt
	deleteDukungan        *sqlx.Stmt
	getPendukungFull      *sqlx.Stmt
	getUsers              *sqlx.Stmt
	changePassword        *sqlx.Stmt
	deleteUser            *sqlx.Stmt
	deleteDukunganByCalon *sqlx.Stmt
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
	r.findIDStmt = r.MustPrepareStmt("SELECT * FROM User WHERE id=?")
	r.findAtDukungan = r.MustPrepareStmt("SELECT * FROM Dukungan WHERE nik =? AND tingkat=?")
	r.findAtPendukung = r.MustPrepareStmt("SELECT * FROM Pendukung WHERE nik=?")
	r.findUsrnameStmt = r.MustPrepareStmt("SELECT * FROM User WHERE username=?")
	r.findPendukungByCalon = r.MustPrepareStmt("select p.id, p.nik, name, phone, witness, gender, status, provinsi, kabupaten, kecamatan, kelurahan, tps, timestamp from Pendukung p inner join Dukungan d on p.nik=d.nik where idCalon=?")
	r.insertUser = r.MustPrepareNamedStmt("INSERT INTO User (id, name, tingkat, username, password, role) VALUES (:id, :name, :tingkat, :username, :password, :role)")
	r.insertDukungan = r.MustPrepareNamedStmt("INSERT INTO Dukungan (id, idCalon, nik, tingkat, status, timestamp) VALUES (:id, :idCalon, :nik, :tingkat, :status, :timestamp)")
	r.insertPendukung = r.MustPrepareNamedStmt("INSERT INTO Pendukung (id, name, nik, provinsi, kabupaten, kecamatan, kelurahan, tps, phone, witness, gender, photo) VALUES (:id, :name, :nik, :provinsi, :kabupaten, :kecamatan, :kelurahan, :tps, :phone, :witness, :gender, :photo)")
	r.confirmDukungan = r.MustPrepareStmt("UPDATE Dukungan SET status=true WHERE nik=? and tingkat=?")
	r.deleteDukungan = r.MustPrepareStmt("DELETE FROM Dukungan WHERE nik=? and tingkat=?")
	r.getPendukungFull = r.MustPrepareStmt("select p.id, p.nik, name, phone, witness, gender, status, provinsi, kabupaten, kecamatan, kelurahan, tps, photo from Pendukung p inner join Dukungan d on p.nik=d.nik where p.nik=? and tingkat=?")
	r.getUsers = r.MustPrepareStmt("select id, name, tingkat from User where id != 1")
	r.changePassword = r.MustPrepareStmt("UPDATE User SET password=? WHERE username=?")
	r.deleteUser = r.MustPrepareStmt("DELETE FROM User WHERE id=?")
	r.deleteDukunganByCalon = r.MustPrepareStmt("DELETE FROM Dukungan WHERE idCalon=?")
	return &r
}

func (db *userRepository) DeleteUser(idCalon string) (success bool, err error) {
	res, errDB := db.deleteUser.Exec(idCalon)
	if errDB != nil {
		err = errDB
		log.Println("Failed to delete User: ", err)
		success = false
	}
	rowUpdated, _ := res.RowsAffected()
	success = (rowUpdated > 0)
	return
}

func (db *userRepository) DeleteDukunganByCalon(idCalon string) (success bool, err error) {
	res, errDB := db.deleteDukunganByCalon.Exec(idCalon)
	if errDB != nil {
		err = errDB
		log.Println("Failed to delete dukungan by calon: ", err)
		success = false
	}
	rowUpdated, _ := res.RowsAffected()
	success = (rowUpdated > 0)
	return
}

func (db *userRepository) ChangePassword(username string, newPassword string) (success bool, err error) {
	res, errDB := db.changePassword.Exec(newPassword, username)
	if errDB != nil {
		err = errDB
		log.Println("Failed to changePassword: ", err)
		success = false
	}
	rowUpdated, _ := res.RowsAffected()
	success = (rowUpdated > 0)
	return
}

func (db *userRepository) GetUsers() (users []UserPart, err error) {
	err = db.getUsers.Select(&users)
	if err != nil {
		log.Println("Error at get all users part,    ", err)
	}
	return
}

func (db *userRepository) ConfirmDukungan(nik string, tingkat string) (success bool, err error) {
	res, errDB := db.confirmDukungan.Exec(nik, tingkat)
	if errDB != nil {
		err = errDB
		log.Println("Failed to confirm dukungan: ", err)
		success = false
	}
	rowUpdated, _ := res.RowsAffected()
	success = (rowUpdated > 0)
	return
}

func (db *userRepository) DeleteDukungan(nik string, tingkat string) (success bool, err error) {
	res, errDB := db.deleteDukungan.Exec(nik, tingkat)
	if errDB != nil {
		err = errDB
		log.Println("Failed to delete dukungan: ", err)
		success = false
	}
	rowUpdated, _ := res.RowsAffected()
	success = (rowUpdated > 0)
	return
}

func (db *userRepository) GetPendukungFull(nik string, tingkat string) (pendukung PendukungFull, err error) {
	err = db.getPendukungFull.Get(&pendukung, nik, tingkat)
	if err != nil {
		log.Println("Error at get Pendukung full,    ", err)
	}
	return
}

func (db *userRepository) FindByID(id string) (usr User, err error) {
	err = db.findIDStmt.Get(&usr, id)
	if err != nil {
		log.Printf("ID: %v , doesn't exist", id)
		log.Println(err)
	}
	return
}

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
