package user_repo

import (
	"context"

	"github.com/Light-ink-yht/admin-hub/internal/domain/user_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
)

var (
	ErrUserDuplicateEmailOrPhone = dao.ErrUserDuplicateEmailOrPhone
)

type UserRepo interface {
	Signup(ctx context.Context, user *user_domain.AA01) error
}

type userRepo struct {
	dao dao.UserDao
}

func NewUserRepo(dao dao.UserDao) UserRepo {
	return &userRepo{
		dao: dao,
	}
}

func (repo *userRepo) Signup(ctx context.Context, user *user_domain.AA01) error {
	return repo.dao.Insert(ctx, user)
}
