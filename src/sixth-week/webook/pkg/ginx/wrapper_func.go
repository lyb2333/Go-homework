package ginx

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Req struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Phone           string `json:"phone"`
	Code            string `json:"code"`
}

type Result struct {
	Code int
	Msg  string
}

var (
	L      *zap.Logger
	vector *prometheus.CounterVec
)

func WrapHandler(bizFn func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("输入错误", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "输入错误"})
			return
		}
		L.Debug("输入参数", zap.Any("req", req))

		res, err := bizFn(ctx, req)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("执行业务逻辑失败", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "执行业务逻辑失败"})
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapClaims[Claims any](
	bizFn func(ctx *gin.Context, uc Claims) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("user")
		if !ok {
			L.Error("未找到用户信息")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
			return
		}
		uc, ok := val.(Claims)
		if !ok {
			L.Error("用户信息类型错误")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "用户信息类型错误"})
			return
		}

		res, err := bizFn(ctx, uc)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("执行业务逻辑失败", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "执行业务逻辑失败"})
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}
