package tcp_proxy_middleware

import (
	"context"
	s "gateway-micro/tcp_server"
	"math"
	"net"
)

const abortIndex int8 = math.MaxInt8 / 2 //最多 63 个中间件

type TcpHandlerFunc func(*TcpSliceRouterContext)

type TcpSliceRouter struct {
	groups []*TcpSliceGroup
}

type TcpSliceGroup struct {
	*TcpSliceRouter
	path     string
	handlers []TcpHandlerFunc
}

type TcpSliceRouterContext struct {
	conn net.Conn
	Ctx  context.Context
	*TcpSliceGroup
	index int8
}

func newTcpSliceRouterContext(conn net.Conn, r *TcpSliceRouter, ctx context.Context) *TcpSliceRouterContext {
	newTcpSliceGroup := &TcpSliceGroup{}
	*newTcpSliceGroup = *r.groups[0] //浅拷贝数组指针,只会使用第一个分组
	c := &TcpSliceRouterContext{conn: conn, TcpSliceGroup: newTcpSliceGroup, Ctx: ctx}
	c.Reset()
	return c
}

func (c *TcpSliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *TcpSliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

type TcpSliceRouterHandler struct {
	coreFunc func(*TcpSliceRouterContext) s.TCPHandler
	router   *TcpSliceRouter
}

func (w *TcpSliceRouterHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	c := newTcpSliceRouterContext(conn, w.router, ctx)
	c.handlers = append(c.handlers, func(c *TcpSliceRouterContext) {
		w.coreFunc(c).ServeTCP(ctx, conn)
	})
	c.Reset()
	c.Next()
}

func NewTcpSliceRouterHandler(coreFunc func(*TcpSliceRouterContext) s.TCPHandler, router *TcpSliceRouter) *TcpSliceRouterHandler {
	return &TcpSliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

func NewTcpSliceRouter() *TcpSliceRouter {
	return &TcpSliceRouter{}
}

func (g *TcpSliceRouter) Group(path string) *TcpSliceGroup {
	if path != "/" {
		panic("only accept path=/")
	}
	return &TcpSliceGroup{
		TcpSliceRouter: g,
		path:           path,
	}
}

func (g *TcpSliceGroup) Use(middlewares ...TcpHandlerFunc) *TcpSliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	for _, oldGroup := range g.TcpSliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.TcpSliceRouter.groups = append(g.TcpSliceRouter.groups, g)
	}
	return g
}

func (c *TcpSliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *TcpSliceRouterContext) Abort() {
	c.index = abortIndex
}

func (c *TcpSliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *TcpSliceRouterContext) Reset() {
	c.index = -1
}
