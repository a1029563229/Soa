package soa

import (
	"lib/soa/sutils"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Server struct {
	middlewares []Middleware
}

var record = make(map[string][]string)

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func routeController() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !sutils.HasKey(record, r.URL.Path) {
				http.NotFound(w, r)
				return
			}
			if !sutils.Includes(record[r.URL.Path], r.Method) {
				http.NotFound(w, r)
				return
			}
			next(w, r)
		}
	}
}

func (s *Server) Use(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *Server) GET(uri string, handle func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	s.setRequest("GET", uri, handle, middlewares...)
}

func (s *Server) POST(uri string, handle func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	s.setRequest("GET", uri, handle, middlewares...)
}

func (s *Server) PUT(uri string, handle func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	s.setRequest("GET", uri, handle, middlewares...)
}

func (s *Server) DELETE(uri string, handle func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	s.setRequest("GET", uri, handle, middlewares...)
}

func (s *Server) setRequest(method string, uri string, handle func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	record[uri] = append(record[uri], method)
	handle = Chain(http.HandlerFunc(handle), routeController())
	for _, m := range s.middlewares {
		handle = Chain(handle, m)
	}
	http.Handle(uri, Chain(handle, middlewares...))
}
