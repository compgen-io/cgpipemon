package model

import (
	"database/sql"
	"github.com/compgen-io/cgpipemon/auth"
)

type User struct {
	Id int
	Username string
	Admin bool
}

func NewUser(db *sql.DB, username string, passwd string, admin bool) (*User, error) {
	enc, err1 := auth.EncryptPass(1, passwd)
	if err1 != nil {
		return nil, err1
	}

	var id int

	err := db.QueryRow("INSERT INTO usr (username, password, is_admin) VALUES ($1, $2, $3) RETURNING id", username, enc, admin).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &User{id, username, admin}, nil
}

func CheckPass(db *sql.DB, username string, plain string) bool {
	var enc string
	err := db.QueryRow("SELECT password FROM usr WHERE username = $1", username).Scan(&enc)
	if err != nil {
		return false
	}
	return auth.CheckPass(plain, enc)
}
