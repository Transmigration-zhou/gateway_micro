package http_proxy_middleware

import (
	"errors"
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			//fmt.Println("item", item)
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			compile, err := regexp.Compile(items[0])
			if err != nil {
				fmt.Println("regexp.Compile err", err)
				continue
			}
			//fmt.Println("before rewrite", c.Request.URL.Path)
			replacePath := compile.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
			//fmt.Println("after rewrite", c.Request.URL.Path)
		}
		c.Next()
	}
}
