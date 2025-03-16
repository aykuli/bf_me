package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/services"
	"bf_me/internal/storage"
	"errors"
	"fmt"
)

var (
	ErrNoMoreUser = errors.New("no more users can register")
)

type SessionsUseCase struct {
	storage *storage.Storage
}

func NewSessionsUseCase(st *storage.Storage) *SessionsUseCase {
	return &SessionsUseCase{storage: st}
}

func (suc *SessionsUseCase) CreateUser(req requests.UserRequestBody) (*models.Session, error) {
	// limiting registered users
	var count int64
	result := suc.storage.DB.Table("users").Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}

	if count > 1 {
		return nil, ErrNoMoreUser
	}

	var u models.User
	var err error
	u.Login = req.Login
	u.PasswordHash, err = services.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password error: %s", err)
	}

	result = suc.storage.DB.Create(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("create user error: %s", result.Error)
	}
	var s models.Session
	s.User = u
	result = suc.storage.DB.Create(&s)
	return &s, result.Error
}

func (suc *SessionsUseCase) Create(req *requests.UserRequestBody) (*models.Session, error) {
	var u models.User
	result := suc.storage.DB.Where("login = ?", req.Login).First(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("no such user error: %s", result.Error)
	}

	wrong := services.WrongPassword(req.Password, u.PasswordHash)
	if wrong {
		return nil, fmt.Errorf("wrong password")
	}
	var sessions []models.Session
	suc.storage.DB.Where("user_id = ?", u.ID).Delete(&sessions)

	var s models.Session
	s.User = u
	result = suc.storage.DB.Create(&s)
	return &s, result.Error
}

func (suc *SessionsUseCase) Delete(sessionId string) error {
	var session models.Session
	result := suc.storage.DB.Where("id = ?", sessionId).First(&session)
	if result.Error != nil {
		return result.Error
	}

	var sessions []models.Session
	result = suc.storage.DB.Where("user_id = ?", session.UserID).Delete(&sessions)
	return result.Error
}

func (suc *SessionsUseCase) Find(sessionId string) (*models.Session, error) {
	var session *models.Session
	result := suc.storage.DB.Where("id = ?", sessionId).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}
