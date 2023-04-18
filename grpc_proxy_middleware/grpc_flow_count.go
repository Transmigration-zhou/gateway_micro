package grpc_proxy_middleware

import (
	"gateway-micro/dao"
	"gateway-micro/public"
	"google.golang.org/grpc"
	"log"
)

// GRPCFlowCountMiddleware 流量统计中间件
func GRPCFlowCountMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			return err
		}
		totalCounter.Increase()
		//dayCount, _ := totalCounter.GetDayData(time.Now())
		//fmt.Printf("totalCounter QPS:%v, dayCount:%v\n", totalCounter.QPS, dayCount)

		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()
		//dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		//fmt.Printf("serviceCounter QPS:%v, dayCount:%v\n", serviceCounter.QPS, dayServiceCount)

		if err := handler(srv, ss); err != nil {
			log.Printf("GRPCFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
