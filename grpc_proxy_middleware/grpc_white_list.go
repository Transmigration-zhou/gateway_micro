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

func GRPCWhiteListMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		addr := peerCtx.Addr.String()
		clientIP := addr[:strings.LastIndex(addr, ":")]

		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, clientIP) {
				return errors.New(fmt.Sprintf("%s not in white ip list", clientIP))
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GRPCWhiteListMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
