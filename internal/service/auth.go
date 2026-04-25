package service

import (
	"context"
	"fmt"

	"github.com/BogdanBratsky/eigth/internal/model"
	"github.com/BogdanBratsky/eigth/pkg/hasher"
)

type AuthRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthService struct {
	repo   AuthRepo
	hasher *hasher.BcryptHasher
}

func NewAuthService(r AuthRepo, h *hasher.BcryptHasher) *AuthService {
	return &AuthService{
		repo:   r,
		hasher: h,
	}
}

func (s *AuthService) Register(ctx context.Context, login, email, password string) error {
	exists, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	if exists != nil {
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
