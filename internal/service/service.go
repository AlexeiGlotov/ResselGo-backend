package service

import (
	"RecurroControl/internal/repository"
	"RecurroControl/models"
)

type Authorization interface {
	CreateUser(user models.SignUpInput) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	CheckAccessKey(key string) (*models.AccessKey, error)
	SetLoginAccessKey(login, key string) error
}

type AccessKeys interface {
	CreateAccessKey(userID int, role string) (string, error)
	GetAccessKey(login, role string) ([]models.AccessKey, error)
}

type Cheats interface {
	GetCheats(role string) ([]models.Cheats, error)
	CreateCheat(cheat *models.Cheats) (int, error)
	UpdateCheat(cheat *models.Cheats) error
}

type LicenseKeys interface {
	CreateLicenseKeys(keys []models.LicenseKeys) error
	GetLicenseKeys(userID, limit, offset int) ([]models.LicenseKeys, error)
}

type Users interface {
	GetUsers(userID int) ([]models.User, error)
	GetUser(userID int) (*models.User, error)
	Ban(userID int) error
	Unban(userID int) error
	Delete(userID int) error
}
type Service struct {
	Authorization
	AccessKeys
	Cheats
	Users
	LicenseKeys
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		AccessKeys:    NewAdmissionService(repos.AccessKeys),
		Cheats:        NewCheatService(repos.Cheats),
		Users:         NewUsersService(repos.Users),
		LicenseKeys:   NewLicenseKeysService(repos.LicenseKeys),
	}
}
