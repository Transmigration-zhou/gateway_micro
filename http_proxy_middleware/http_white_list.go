package http_proxy_middleware

import (
	"errors"
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"strings"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
