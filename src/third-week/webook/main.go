package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"homework/src/third-week/webook/config"
	"homework/src/third-week/webook/internal/repository"
	"homework/src/third-week/webook/internal/repository/dao"
	"homework/src/third-week/webook/internal/service"
	"homework/src/third-week/webook/internal/web"
	"homework/src/third-week/webook/internal/web/middleware"
	"homework/src/third-week/webook/pkg/ginx/middleware/ratelimit"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()

	server := initWebServer()
	initUserHdl(db, server)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello，启动成功了！")
	})
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,

		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 这个是允许前端访问你的后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
		//AllowHeaders: []string{"content-type"},
		//AllowMethods: []string{"POST"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这是我的 Middleware")
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	server.Use(ratelimit.NewBuilder(redisClient,
		time.Second, 1).Build())

	useJWT(server)
	//useSession(server)
	return server
}

func useJWT(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	store := cookie.NewStore([]byte("secret"))
	// 基于内存的实现
	//store := memstore.NewStore([]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
	//	[]byte("eF1`yQ9>yT1`tH1,sJ0.zD8;mZ9~nC6("))
	//store, err := redis.NewStore(16, "tcp",
	//	"localhost:6379", "",
	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
