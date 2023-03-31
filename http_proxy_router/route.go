package http_proxy_router

import (
	"gateway-micro/controller"
	"gateway-micro/http_proxy_middleware"
	"gateway-micro/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	//优化点：gin.Default()使用Logger(), Recovery()中间件，会打印请求日志，消耗性能IO
	//router := gin.Default()
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oAuth := router.Group("/oauth")
	oAuth.Use(
		middleware.TranslationMiddleware(),
	)
	{
		controller.OAuthRegister(oAuth)
	}

	router.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(),
		http_proxy_middleware.HTTPFlowCountMiddleware(),
		http_proxy_middleware.HTTPFlowLimitMiddleware(),
		http_proxy_middleware.HTTPJwtAuthTokenMiddleware(),
		http_proxy_middleware.HTTPJwtFlowCountMiddleware(),
		http_proxy_middleware.HTTPWhiteListMiddleware(),
		http_proxy_middleware.HTTPBlackListMiddleware(),
		http_proxy_middleware.HTTPHeaderTransferMiddleware(),
		http_proxy_middleware.HTTPStripUriMiddleware(),
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),
		http_proxy_middleware.HTTPReverseProxyMiddleware(),
	)

	return router
}
