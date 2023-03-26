package http_proxy_middleware

import (
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"github.com/gin-gonic/gin"
)

// HTTPAccessModeMiddleware 基于请求信息，匹配接入方式中间件
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		//fmt.Println("matched service", public.Object2Json(service))
		c.Set("service", service)
		c.Next()
	}
}
