package service

import (
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/token"
	"auth-service/internal/utils"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type AuthService struct {
	userRepo    *repository.UserRepository    // ✅ lowercase
	refreshRepo *repository.RefreshRepository // ✅ lowercase
	config      *config.Config                // ✅ lowercase
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshRepo *repository.RefreshRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		config:      cfg,
	}
}

// jwtExpire แปลง JWTExpireHour จาก config → time.Duration
// ถ้าค่าว่างหรือ parse ไม่ได้ fallback เป็น 24h พร้อม log เตือน
func (s *AuthService) jwtExpire() time.Duration {
	hours, err := strconv.Atoi(s.config.JWTExpireHour)
	if err != nil || hours <= 0 {
		return 24 * time.Hour // fallback
	}
	return time.Duration(hours) * time.Hour
}

func (s *AuthService) Register(req dto.RegisterRequest) error {
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
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

	return s.userRepo.Create(&user)
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := token.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		s.config.JWTSecret,
		s.jwtExpire(), // ✅ ใช้ expire จาก config
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

	if err := s.refreshRepo.Create(&refreshModel); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(req dto.RefreshRequest) (*dto.LoginResponse, error) {
	refresh, err := s.refreshRepo.FindByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refresh.Revoked {
		return nil, errors.New("refresh token revoked")
	}

	if time.Now().After(refresh.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	if err := s.refreshRepo.Revoke(refresh.ID.String()); err != nil {
		return nil, err
	}

	// ✅ ใช้ FindByID แทนการเรียก .DB โดยตรง
	user, err := s.userRepo.FindByID(refresh.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	newAccessToken, err := token.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		s.config.JWTSecret,
		s.jwtExpire(), // ✅ ใช้ expire จาก config
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

	if err := s.refreshRepo.Create(&refreshModel); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Username:     user.Username,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// ✅ Logout อยู่ใน service layer — handler ไม่ต้องรู้จัก repo เลย
func (s *AuthService) Logout(req dto.RefreshRequest) error {
	refresh, err := s.refreshRepo.FindByToken(req.RefreshToken)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	if refresh.Revoked {
		return errors.New("token already revoked")
	}

	return s.refreshRepo.Revoke(refresh.ID.String())
}
