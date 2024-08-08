package repository

import (
	"context"
	"homework/src/second-week/webook/internal/domain"
	"homework/src/second-week/webook/internal/repository/dao"
	"strconv"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) Update(ctx context.Context, u domain.User) error {
	int64Id, err := strconv.ParseInt(u.Id, 10, 64)
	if err != nil {
		return err
	}
	return repo.dao.Update(ctx, dao.User{
		Id:       int64Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Aboutme:  u.Aboutme,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx context.Context, id string) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       strconv.FormatInt(u.Id, 10),
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Phone:    u.Phone,
		Aboutme:  u.Aboutme,
		Birthday: u.Birthday,
	}
}
