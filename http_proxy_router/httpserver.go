package http_proxy_router

import (
	"context"
	"gateway-micro/cert_file"
	"gateway-micro/common/lib"
	"gateway-micro/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpProxyRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}

	log.Printf("[INFO] HttpProxyRun %s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[ERROR] HttpProxyRun %s, err: %v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpProxyStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] HttpProxyStop %s, err: %v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
	log.Printf("[INFO] HttpProxyStop %s stopped\n", lib.GetStringConf("proxy.http.addr"))
}

func HttpsProxyRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}

	log.Printf("[INFO] HttpsProxyRun %s\n", lib.GetStringConf("proxy.https.addr"))
	if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[ERROR] HttpsProxyRun %s, err: %v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpsProxyStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] HttpsProxyStop %s, err: %v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
	log.Printf("[INFO] HttpsProxyStop %s stopped\n", lib.GetStringConf("proxy.https.addr"))
}
