//go:build wireinject

package wire

import (
	"homework/src/fourth-week/wire/repository"
	"homework/src/fourth-week/wire/repository/dao"

	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, InitDB, dao.NewUserDAO)
	return &repository.UserRepository{}
}
