package http_proxy_middleware

import (
	"errors"
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"github.com/gin-gonic/gin"
	"strings"
)

// HTTPHeaderTransferMiddleware header头转换中间件
func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfer, ",") {
			items := strings.Split(item, " ")
			if len(items) == 3 && (items[0] == "add" || items[0] == "edit") {
				c.Request.Header.Set(items[1], items[2])
			} else if len(items) == 2 && items[0] == "del" {
				c.Request.Header.Del(items[1])
			}
		}
		c.Next()
	}
}
