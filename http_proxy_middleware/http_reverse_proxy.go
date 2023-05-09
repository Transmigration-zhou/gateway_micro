package http_proxy_middleware

import (
	"gateway-micro/dao"
	"gateway-micro/middleware"
	"gateway-micro/reverse_proxy"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HTTPReverseProxyMiddleware 反向代理中间件，反向代理后才能访问到下游服务器
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 负载均衡
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}

		// 连接池
		trans, err := dao.TransporterHandler.GetTransporter(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}

		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request) // 执行下游服务器
		c.Abort()
	}
}
