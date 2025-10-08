package user_repo

import (
	"context"

	"github.com/Light-ink-yht/admin-hub/internal/domain/user_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
)

var (
	ErrUserDuplicateEmailOrPhone = dao.ErrUserDuplicateEmailOrPhone
	ErrUserNotFound              = dao.ErrUserNotFound
)

type UserRepo interface {
	Signup(ctx context.Context, user *user_domain.AA01) error
	FindByEmail(ctx context.Context, aaa002 string) (*user_domain.AA01, error)
	FindByPhone(ctx context.Context, aaa003 string) (*user_domain.AA01, error)
	FindByUserId(ctx context.Context, aaa001 int64) (*user_domain.AA01, error)
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

func (repo *userRepo) FindByEmail(ctx context.Context, aaa002 string) (*user_domain.AA01, error) {
	return repo.dao.FindByEmail(ctx, aaa002)
}

func (repo *userRepo) FindByPhone(ctx context.Context, aaa003 string) (*user_domain.AA01, error) {
	return repo.dao.FindByPhone(ctx, aaa003)
}

func (repo *userRepo) FindByUserId(ctx context.Context, aaa001 int64) (*user_domain.AA01, error) {
	return repo.dao.FindByUserId(ctx, aaa001)
}
