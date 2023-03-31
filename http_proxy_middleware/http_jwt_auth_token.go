package http_proxy_middleware

import (
	"errors"
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"strings"
)

func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// decode jwt token
		// tenant_id 与 tenant_list 取得 tenant
		// tenant 放到 gin.context
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		tenantMatched := false
		if token != "" {
			claims, err := public.JwtDecode(token)
			if err != nil {
				middleware.ResponseError(c, 2002, err)
				c.Abort()
				return
			}
			tenantList := dao.TenantManagerHandler.GetTenantList()
			for _, tenant := range tenantList {
				if tenant.TenantID == claims.Issuer {
					c.Set("tenant", tenant)
					tenantMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !tenantMatched {
			middleware.ResponseError(c, 2003, errors.New("not match valid tenant"))
			c.Abort()
			return
		}
		c.Next()
	}
}
