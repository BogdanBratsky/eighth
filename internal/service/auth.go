package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BogdanBratsky/eigth/internal/model"
	"github.com/BogdanBratsky/eigth/pkg/hasher"
	"github.com/BogdanBratsky/eigth/pkg/token"
)

type AuthRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type AuthService struct {
	repo   AuthRepo
	hasher *hasher.BcryptHasher
	jwt    *token.JWTManager
	logger *slog.Logger
}

func NewAuthService(
	r AuthRepo,
	h *hasher.BcryptHasher,
	j *token.JWTManager,
	l *slog.Logger,
) *AuthService {
	return &AuthService{
		repo:   r,
		hasher: h,
		jwt:    j,
		logger: l,
	}
}

func (s *AuthService) Register(ctx context.Context, login, email, password string) error {
	// check email
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		s.logger.Warn("register attempt with existing email",
			"email", email,
		)
		return ErrUserExists
	}

	// check login
	_, err = s.repo.GetByLogin(ctx, login)
	if err == nil {
		s.logger.Warn("register attempt with existing login",
			"login", login,
		)
		return ErrUserExists
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		s.logger.Error("password hashing failed",
			"error", err,
			"login", login,
			"email", email,
		)
		return fmt.Errorf("hash error: %w", err)
	}

	user := &model.User{
		Login:        login,
		Email:        email,
		PasswordHash: hash,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user",
			"error", err,
			"email", email,
			"login", login,
		)
		return fmt.Errorf("create error: %w", err)
	}

	s.logger.Info("user registered successfully",
		"email", email,
		"login", login,
	)

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("login failed: user not found",
			"email", email,
		)
		return "", ErrUserNotFound
	}

	err = s.hasher.Compare(user.PasswordHash, password)
	if err != nil {
		s.logger.Warn("login failed: invalid password",
			"email", email,
		)
		return "", ErrInvalidPassword
	}

	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error("failed to generate jwt token",
			"error", err,
			"user_id", user.ID,
		)
		return "", err
	}

	s.logger.Info("user logged in successfully",
		"user_id", user.ID,
		"email", email,
	)

	return token, nil
}
