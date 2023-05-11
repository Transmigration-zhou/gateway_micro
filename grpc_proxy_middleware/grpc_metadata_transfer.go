package grpc_proxy_middleware

import (
	"gateway-micro/dao"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

// GRPCMetadataTransferMiddleware metadata转换中间件
func GRPCMetadataTransferMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		for _, item := range strings.Split(serviceDetail.GRPCRule.HeaderTransfer, ",") {
			items := strings.Split(item, " ")
			if len(items) == 3 && (items[0] == "add" || items[0] == "edit") {
				md.Set(items[1], items[2])
			} else if len(items) == 2 && items[0] == "del" {
				delete(md, items[1])
			}
		}

		if err := ss.SetHeader(md); err != nil {
			return errors.WithMessage(err, "SetHeader")
		}

		if err := handler(srv, &wrappedStream{ss, metadata.NewIncomingContext(ss.Context(), md)}); err != nil {
			log.Printf("GRPCMetadataTransferMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
