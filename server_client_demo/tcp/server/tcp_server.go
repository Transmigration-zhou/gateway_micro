package server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed     = errors.New("tcp: Server closed")
	ErrAbortHandler     = errors.New("tcp: abort TCPHandler")
	ServerContextKey    = &contextKey{"tcp-server"}
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() {
	oc.closeErr = oc.Listener.Close()
}

type TCPHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

type TcpServer struct {
	Addr    string
	Handler TCPHandler
	err     error
	BaseCtx context.Context

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	KeepAliveTimeout time.Duration

	mu         sync.Mutex
	inShutdown int32
	doneChan   chan struct{}
	l          *onceCloseListener
}

func (s *TcpServer) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutdown) != 0
}

func (s *TcpServer) ListenAndServe() error {
	if s.shuttingDown() {
		return ErrServerClosed
	}
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	addr := s.Addr
	if addr == "" {
		return errors.New("need addr")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(tcpKeepAliveListener{
		ln.(*net.TCPListener)})
}

func (s *TcpServer) Close() error {
	atomic.StoreInt32(&s.inShutdown, 1)
	close(s.doneChan) //关闭channel
	s.l.Close()       //执行listener关闭
	return nil
}

func (s *TcpServer) Serve(l net.Listener) error {
	s.l = &onceCloseListener{Listener: l}
	defer s.l.Close() //执行listener关闭
	if s.BaseCtx == nil {
		s.BaseCtx = context.Background()
	}
	baseCtx := s.BaseCtx
	ctx := context.WithValue(baseCtx, ServerContextKey, s)
	for {
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-s.getDoneChan():
				return ErrServerClosed
			default:
			}
			fmt.Printf("accept fail, err: %v\n", e)
			continue
		}
		c := s.newConn(rw)
		go c.serve(ctx)
	}
	return nil
}

func (s *TcpServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: s,
		rwc:    rwc,
	}
	// 设置参数
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := c.server.KeepAliveTimeout; d != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
	return c
}

func (s *TcpServer) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func ListenAndServe(addr string, handler TCPHandler) error {
	server := &TcpServer{Addr: addr, Handler: handler, doneChan: make(chan struct{})}
	return server.ListenAndServe()
}
