package grpc_proxy_middleware

import (
	"gateway-micro/dao"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

// GRPCHeaderTransferMiddleware header头转换中间件
func GRPCHeaderTransferMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfer, ",") {
			items := strings.Split(item, " ")
			if len(items) == 3 && (items[0] == "add" || items[0] == "edit") {
				md.Set(items[1], items[2])
			} else if len(items) == 2 && items[0] == "del" {
				delete(md, items[1])
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GRPCHeaderTransferMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
