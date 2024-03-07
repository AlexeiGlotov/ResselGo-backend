package repository

import (
	"database/sql"

	todo "RecurroControl"
)

const (
	usersTable     = "users"
	admissionTable = "reg_admission"
	cheatTable     = "cheats"
)

type Authorization interface {
	CreateUser(user todo.SignUpInput) (int, error)
	GetUser(username, password string) (todo.User, error)
	CheckKeyAdmission(key string) (string, error)
	SetLoginAdmission(login, key string) error
}

type Admission interface {
	CreateKey(userID int) (string, error)
	GetKey() ([]todo.RegAdmission, error)
}

type Cheat interface {
	GetCheats() ([]todo.StCheats, error)
}

type Repository struct {
	Authorization
	Admission
	Cheat
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(db),
		Admission:     NewAdmissionSql(db),
		Cheat:         NewCheatSql(db),
	}
}
