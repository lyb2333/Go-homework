//go:build wireinject

package main

import (
	"homework/src/fourth-week/webook/internal/repository"
	"homework/src/fourth-week/webook/internal/repository/cache"
	"homework/src/fourth-week/webook/internal/repository/dao"
	"homework/src/fourth-week/webook/internal/service"
	"homework/src/fourth-week/webook/internal/web"
	"homework/src/fourth-week/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		// DAO 部分
		dao.NewUserDAO,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
