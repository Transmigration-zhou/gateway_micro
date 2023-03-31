package http_proxy_middleware

import (
	"errors"
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantInterface, ok := c.Get("tenant")
		if !ok {
			c.Next()
			return
		}
		tenant := tenantInterface.(*dao.Tenant)

		tenantCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTenantPrefix + tenant.TenantID)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		tenantCounter.Increase()
		if tenant.Qpd > 0 && tenantCounter.TotalCount > tenant.Qpd {
			middleware.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", tenant.Qpd, tenantCounter.TotalCount)))
			c.Abort()
			return
		}
		c.Next()
	}
}
