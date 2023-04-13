package http_proxy_middleware

import (
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		var whileIpList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		var blackIpList []string
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if public.InStringSlice(blackIpList, c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s in black ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
