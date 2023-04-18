package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"gateway-micro/dao"
	"gateway-micro/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func GRPCJwtFlowCountMiddleware() grpc.StreamServerInterceptor {
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

			tenantCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTenantPrefix + tenant.TenantID)
			if err != nil {
				return err
			}
			tenantCounter.Increase()
			if tenant.Qpd > 0 && tenantCounter.TotalCount > tenant.Qpd {
				return errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", tenant.Qpd, tenantCounter.TotalCount))
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GRPCJwtFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
