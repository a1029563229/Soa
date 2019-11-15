package soa

import (
	"fmt"
	"lib/soa/sutils"
	"net/http"
	"strconv"
)

type Server struct {
}

type Request struct {
	URL    string
	Method string
}
type Ctx struct {
	w       http.ResponseWriter
	r       *http.Request
	Request Request
}

func (ctx *Ctx) init() {
	ctx.Request.URL = ctx.r.URL.Path
	ctx.Request.Method = ctx.r.Method
}

func (ctx *Ctx) End(status int, message string) {
	ctx.w.WriteHeader(status)
	fmt.Fprintf(ctx.w, message)
}

type Handle func(ctx *Ctx)

type Middleware func(Handle) Handle

var routes = make(map[string][]string)

func routeController() Middleware {
	return func(next Handle) Handle {
		return func(ctx *Ctx) {
			if !sutils.Includes(routes[ctx.Request.URL], ctx.Request.Method) {
				http.NotFound(ctx.w, ctx.r)
				return
			}
			next(ctx)
		}
	}
}

func (s *Server) Listen(port int) {
	p := strconv.Itoa(port)
	fmt.Println("server is listening in http://localhost:" + p)
	http.ListenAndServe(":"+p, nil)
}

func (s *Server) GET(uri string, handle Handle, middlewares ...Middleware) {
	s.SetRequest("GET", uri, handle, middlewares...)
}

func (s *Server) SetRequest(method string, uri string, handle Handle, middlewares ...Middleware) {
	routes[uri] = append(routes[uri], method)
	handle = Chain(handle, routeController())
	handle = Chain(handle, middlewares...)
	http.HandleFunc(uri, http.HandlerFunc(ctxInject(handle)))
}

func ctxInject(handle Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Ctx{w: w, r: r}
		ctx.init()
		handle(&ctx)
	}
}

func Chain(h Handle, middlewares ...Middleware) Handle {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
