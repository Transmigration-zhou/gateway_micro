package main

import (
	"gateway-micro/common/lib"
	"gateway-micro/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//如果configPath为空，则从命令行 '-config=./conf/prod/' 中读取
	lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
	defer lib.Destroy()
	router.HttpServerRun() // router初始化和服务开启

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
