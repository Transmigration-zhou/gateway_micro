package grpc_proxy_router

import (
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/grpc_proxy_middleware"
	"gateway-micro/public"
	"gateway-micro/reverse_proxy"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"log"
	"net"
)

var grpcServerList []*GrpcServer

type GrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcProxyRun() {
	serviceList := dao.ServiceManagerHandler.GetServiceList(public.LoadTypeGRPC)
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)

			lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf("[ERROR] GetGrpcLoadBalancer %s, err: %v\n", addr, err)
				return
			}

			listener, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("[ERROR] GrpcListen %s, err: %v\n", addr, err)
			}

			grpcHandler := reverse_proxy.NewGrpcLoadBalanceHandler(lb)
			server := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpc_proxy_middleware.GRPCFlowCountMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCFlowLimitMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCJwtAuthTokenMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCJwtFlowCountMiddleware(),
					grpc_proxy_middleware.GRPCJwtFlowLimitMiddleware(),
					grpc_proxy_middleware.GRPCWhiteListMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCBlackListMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCHeaderTransferMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler),
			)

			grpcServer := &GrpcServer{
				Addr:   addr,
				Server: server,
			}
			grpcServerList = append(grpcServerList, grpcServer)
			log.Printf("[INFO] GrpcProxyRun %s\n", addr)
			if err := server.Serve(listener); err != nil {
				log.Fatalf("[ERROR] GrpcProxyRun %s, err: %v\n", addr, err)
			}
		}(tempItem)
	}
}

func GrpcProxyStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf("[INFO] GrpcProxyStop %s stopped\n", grpcServer.Addr)
	}
}
