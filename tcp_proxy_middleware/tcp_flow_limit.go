package tcp_proxy_middleware

import (
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/public"
	"strings"
)

// TCPFlowLimitMiddleware 限流中间件
func TCPFlowLimitMiddleware() TcpHandlerFunc {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit),
			)
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("%s service flow limit %v", serviceDetail.Info.ServiceName, serviceDetail.AccessControl.ServiceFlowLimit)))
				c.Abort()
				return
			}
		}

		split := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(split) == 2 {
			clientIP = split[0]
		}

		if serviceDetail.AccessControl.ClientIpFlowLimit > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIpFlowLimit),
			)
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("%v client ip flow limit %v", clientIP, serviceDetail.AccessControl.ClientIpFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
