package tcp_proxy_middleware

import (
	"gateway-micro/dao"
	"gateway-micro/public"
)

// TCPFlowCountMiddleware 流量统计中间件
func TCPFlowCountMiddleware() TcpHandlerFunc {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()
		//dayCount, _ := totalCounter.GetDayData(time.Now())
		//fmt.Printf("totalCounter QPS:%v, dayCount:%v\n", totalCounter.QPS, dayCount)

		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		//dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		//fmt.Printf("serviceCounter QPS:%v, dayCount:%v\n", serviceCounter.QPS, dayServiceCount)

		c.Next()
	}
}
