// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"homework/src/fifth-week/webook/internal/repository"
	"homework/src/fifth-week/webook/internal/repository/cache"
	"homework/src/fifth-week/webook/internal/repository/dao"
	"homework/src/fifth-week/webook/internal/service"
	"homework/src/fifth-week/webook/internal/web"
	"homework/src/fifth-week/webook/ioc"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.InitGinMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache.(*cache.InMemoryCodeCache))
	smsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService, service.Config{})
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}
