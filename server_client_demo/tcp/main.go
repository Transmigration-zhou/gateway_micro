package main

import (
	"context"
	"gateway-micro/tcp_server"
	"log"
	"net"
)

var (
	addr = ":8042"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func main() {
	server := tcp_server.TcpServer{
		Addr:    addr,
		Handler: &tcpHandler{},
	}
	log.Println("Starting tcpserver at " + addr)
	server.ListenAndServe()
}
