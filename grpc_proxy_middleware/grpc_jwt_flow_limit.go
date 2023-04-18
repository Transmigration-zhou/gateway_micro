package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

func GRPCJwtFlowLimitMiddleware() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		tenants := md.Get("tenant")

		if len(tenants) > 0 {
			tenant := &dao.Tenant{}
			if err := json.Unmarshal([]byte(tenants[0]), tenant); err != nil {
				return err
			}

			peerCtx, ok := peer.FromContext(ss.Context())
			if !ok {
				return errors.New("peer not found with context")
			}
			addr := peerCtx.Addr.String()
			clientIP := addr[:strings.LastIndex(addr, ":")]

			if tenant.Qps > 0 {
				clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
					public.FlowTenantPrefix+tenant.TenantID+"_"+clientIP,
					float64(tenant.Qps),
				)
				if err != nil {
					return err
				}
				if !clientLimiter.Allow() {
					return errors.New(fmt.Sprintf("%v flow limit %v", clientIP, tenant.Qps))
				}
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GRPCJwtFlowLimitMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
