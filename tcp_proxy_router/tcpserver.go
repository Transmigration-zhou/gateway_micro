package tcp_proxy_router

import (
	"context"
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/public"
	"gateway-micro/reverse_proxy"
	"gateway-micro/tcp_proxy_middleware"
	"gateway-micro/tcp_server"
	"log"
)

var tcpServerList []*tcp_server.TcpServer

func TcpProxyRun() {
	serviceList := dao.ServiceManagerHandler.GetServiceList(public.LoadTypeTCP)
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.TCPRule.Port)

			lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf("[INFO] GetTcpLoadBalancer %s, err: %v\n", addr, err)
				return
			}

			//构建路由及设置中间件
			router := tcp_proxy_middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				tcp_proxy_middleware.TCPFlowCountMiddleware(),
				tcp_proxy_middleware.TCPFlowLimitMiddleware(),
				tcp_proxy_middleware.TCPWhiteListMiddleware(),
				tcp_proxy_middleware.TCPBlackListMiddleware(),
			)

			//构建回调handler
			routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
				func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
					return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, lb)
				},
				router,
			)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)

			tcpServer := &tcp_server.TcpServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf("[INFO] TcpProxyRun %s\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Fatalf("[ERROR] TcpProxyRun %s, err: %v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpProxyStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf("[INFO] TcpProxyStop %s stopped\n", tcpServer.Addr)
	}
}
