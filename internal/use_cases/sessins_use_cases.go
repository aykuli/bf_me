package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/services"
	"bf_me/internal/storage"
	"fmt"
)

type SessionsUseCase struct {
	storage *storage.Storage
}

func NewSessionsUseCase(st *storage.Storage) *SessionsUseCase {
	return &SessionsUseCase{storage: st}
}

func (euc *SessionsUseCase) CreateUser(req *requests.UserRequestBody) (*models.Session, error) {
	var u *models.User
	var err error
	u.Login = req.Login
	u.PasswordHash, err = services.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password error: %s", err)
	}
	result := euc.storage.DB.Create(u)
	if result.Error != nil {
		return nil, fmt.Errorf("create user error: %s", result.Error)
	}
	var s *models.Session
	s.User = *u
	result = euc.storage.DB.Create(s)
	return s, result.Error
}

func (euc *SessionsUseCase) Create(req *requests.UserRequestBody) (*models.Session, error) {
	var u *models.User
	result := euc.storage.DB.Where("login = ?", req.Login).First(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("no such user error: %s", result.Error)
	}

	passwordHash, err := services.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	wrong := services.WrongPassword(passwordHash, u.PasswordHash)
	if wrong {
		return nil, fmt.Errorf("wrong password")
	}
	var sessions []models.Session
	euc.storage.DB.Where("user_id = ?", u.ID).Delete(&sessions)

	var s *models.Session
	s.User = *u
	result = euc.storage.DB.Create(s)
	return s, result.Error
}

func (euc *SessionsUseCase) Delete(sessionId string) error {
	var session models.Session
	result := euc.storage.DB.First(&session, sessionId)
	if result.Error != nil {
		return result.Error
	}

	var sessions []models.Session
	result = euc.storage.DB.Where("user_id = ?", session.User.ID).Delete(&sessions)
	return result.Error
}
