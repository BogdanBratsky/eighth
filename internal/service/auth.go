package service

import (
	"context"
	"fmt"
	"log"
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
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return ErrUserExists
	}

	_, err = s.repo.GetByLogin(ctx, login)
	if err == nil {
		return ErrUserExists
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("hash error: %w", err)
	}

	user := &model.User{
		Login:        login,
		Email:        email,
		PasswordHash: hash,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Println("repo.GetByEmail error")
		return "", ErrUserNotFound
	}

	err = s.hasher.Compare(user.PasswordHash, password)
	if err != nil {
		log.Println("hasher.Compare error")
		return "", ErrInvalidPassword
	}

	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		log.Println("jwt.GenerateToken error")
		return "", err
	}

	return token, nil
}
