package service

import (
	"context"

	"github.com/BogdanBratsky/eigth/internal/model"
	"github.com/BogdanBratsky/eigth/pkg/hasher"
)

type UserRepo interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]model.User, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo, h *hasher.BcryptHasher) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Create(ctx context.Context)
