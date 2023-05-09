package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	Addr string
}

func (r *HttpServer) Run() {
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	mux.HandleFunc("/httphttp/timeout", r.TimeoutHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

func (r *HttpServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Http-Ip=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Http-Ip"))
	header := fmt.Sprintf("headers =%v\n", req.Header)
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
	io.WriteString(w, header)
}

func (r *HttpServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (r *HttpServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(2 * time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}

func main() {
	rs1 := &HttpServer{Addr: "127.0.0.1:1234"}
	rs1.Run()
	rs2 := &HttpServer{Addr: "127.0.0.1:4321"}
	rs2.Run()

	//监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
