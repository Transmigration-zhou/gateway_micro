package grpc_proxy_middleware

import (
	"gateway-micro/dao"
	"gateway-micro/public"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

func GRPCJwtAuthTokenMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		auth := ""
		auths := md.Get("Authorization")
		if len(auths) > 0 {
			auth = auths[0]
		}
		token := strings.ReplaceAll(auth, "Bearer ", "")
		tenantMatched := false
		if token != "" {
			claims, err := public.JwtDecode(token)
			if err != nil {
				return err
			}
			tenantList := dao.TenantManagerHandler.GetTenantList()
			for _, tenant := range tenantList {
				if tenant.TenantID == claims.Issuer {
					md.Set("tenant", public.Object2Json(tenant))
					tenantMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !tenantMatched {
			return errors.New("not match valid tenant")
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GRPCJwtAuthTokenMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
