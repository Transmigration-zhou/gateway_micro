package http_proxy_middleware

import (
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		if serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri == 1 {
			//fmt.Println("替换前Path", c.Request.URL.Path)
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)
			//fmt.Println("替换后Path", c.Request.URL.Path)
		}
		c.Next()
	}
}
