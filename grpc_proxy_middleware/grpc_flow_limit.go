package grpc_proxy_middleware

import (
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

// GRPCFlowLimitMiddleware 限流中间件
func GRPCFlowLimitMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit),
			)
			if err != nil {
				return err
			}
			if !serviceLimiter.Allow() {
				return errors.New(fmt.Sprintf("%s service flow limit %v", serviceDetail.Info.ServiceName, serviceDetail.AccessControl.ServiceFlowLimit))
			}
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		addr := peerCtx.Addr.String()
		clientIP := addr[:strings.LastIndex(addr, ":")]

		if serviceDetail.AccessControl.ClientIpFlowLimit > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIpFlowLimit),
			)
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return errors.New(fmt.Sprintf("%v client ip flow limit %v", clientIP, serviceDetail.AccessControl.ClientIpFlowLimit))
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GRPCFlowLimitMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
