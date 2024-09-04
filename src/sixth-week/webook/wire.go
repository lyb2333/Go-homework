//go:build wireinject

package main

import (
	"homework/src/sixth-week/webook/internal/repository"
	"homework/src/sixth-week/webook/internal/repository/cache"
	"homework/src/sixth-week/webook/internal/repository/dao"
	"homework/src/sixth-week/webook/internal/service"
	"homework/src/sixth-week/webook/internal/web"
	ijwt "homework/src/sixth-week/webook/internal/web/jwt"
	"homework/src/sixth-week/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		// DAO 部分
		dao.NewUserDAO,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
