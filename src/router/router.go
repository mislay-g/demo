package router

import (
	"demo/src/config"
	"demo/src/logger"
	"github.com/gin-gonic/gin"
)

func New(c *config.ServerConfig) *gin.Engine {
	gin.SetMode(c.Mode)

	r := gin.New()
	// 注册zap相关中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	configRoute(r)
	return r
}

func configRoute(r *gin.Engine) {
	service := r.Group("/v1")
	{
		service.GET("/kline", getKLine)
	}

}
