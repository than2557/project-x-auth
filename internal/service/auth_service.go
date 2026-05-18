package service

import (
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/token"
	"auth-service/internal/utils"
	"errors"
	"time"
)

type AuthService struct {
	UserRepo    *repository.UserRepository
	RefreshRepo *repository.RefreshRepository
	Config      *config.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshRepo *repository.RefreshRepository,
	cfg *config.Config,
) *AuthService {

	return &AuthService{
		UserRepo:    userRepo,
		RefreshRepo: refreshRepo,
		Config:      cfg,
	}
}

func (s *AuthService) Register(
	req dto.RegisterRequest,
) error {

	existingUser, _ := s.UserRepo.FindByEmail(req.Email)

	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return err
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	return s.UserRepo.Create(&user)
}

func (s *AuthService) Login(
	req dto.LoginRequest,
) (*dto.LoginResponse, error) {

	user, err := s.UserRepo.FindByEmail(req.Email)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	match := utils.CheckPassword(
		req.Password,
		user.Password,
	)

	if !match {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := token.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		s.Config.JWTSecret,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()

	if err != nil {
		return nil, err
	}

	refreshModel := model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = s.RefreshRepo.Create(&refreshModel)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(
	req dto.RefreshRequest,
) (*dto.LoginResponse, error) {

	refresh, err := s.RefreshRepo.
		FindByToken(req.RefreshToken)

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refresh.Revoked {
		return nil, errors.New("refresh token revoked")
	}

	if time.Now().After(refresh.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	err = s.RefreshRepo.Revoke(refresh.ID.String())

	if err != nil {
		return nil, err
	}

	var user model.User

	err = s.UserRepo.DB.
		First(&user, "id = ?", refresh.UserID).
		Error

	if err != nil {
		return nil, err
	}

	newAccessToken, err := token.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		s.Config.JWTSecret,
	)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := token.GenerateRefreshToken()

	if err != nil {
		return nil, err
	}

	refreshModel := model.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = s.RefreshRepo.Create(&refreshModel)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Username:     user.Username,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
