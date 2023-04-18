package main

import (
	"flag"
	"gateway-micro/common/lib"
	"gateway-micro/dao"
	"gateway-micro/grpc_proxy_router"
	"gateway-micro/http_proxy_router"
	"gateway-micro/router"
	"gateway-micro/tcp_proxy_router"
	"os"
	"os/signal"
	"syscall"
)

var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/, default: ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *endpoint == "dashboard" {
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		router.HttpServerRun() // router初始化和服务开启

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		dao.ServiceManagerHandler.LoadOnce()
		dao.TenantManagerHandler.LoadOnce()

		go func() {
			tcp_proxy_router.TcpProxyRun()
		}()
		go func() {
			grpc_proxy_router.GrpcProxyRun()
		}()
		go func() {
			http_proxy_router.HttpProxyRun()
		}()
		go func() {
			http_proxy_router.HttpsProxyRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		tcp_proxy_router.TcpProxyStop()
		grpc_proxy_router.GrpcProxyStop()
		http_proxy_router.HttpProxyStop()
		http_proxy_router.HttpsProxyStop()
	}
}
