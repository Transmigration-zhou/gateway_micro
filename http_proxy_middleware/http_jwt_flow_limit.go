package http_proxy_middleware

import (
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantInterface, ok := c.Get("tenant")
		if !ok {
			c.Next()
			return
		}
		tenant := tenantInterface.(*dao.Tenant)

		if tenant.Qps > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowTenantPrefix+tenant.TenantID+"_"+c.ClientIP(),
				float64(tenant.Qps),
			)
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), tenant.Qps)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
