package tcp_proxy_middleware

import (
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/public"
	"strings"
)

func TCPWhiteListMiddleware() TcpHandlerFunc {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		split := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(split) == 2 {
			clientIP = split[0]
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, clientIP) {
				c.conn.Write([]byte(fmt.Sprintf("%s not in white ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
