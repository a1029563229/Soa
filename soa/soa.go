package soa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	ctx.w.Header().Set("content-type", "application/json")
	ctx.w.WriteHeader(status)
	fmt.Fprintf(ctx.w, message)
}

func (ctx *Ctx) JSON(jsonData interface{}) string {
	jsonBytes, err := json.Marshal(&jsonData) //转换成 JSON，返回的是 byte[]
	if err != nil {
		log.Println(err)
		return "Parse JSON Has Error"
	}

	return string(jsonBytes)
}

func (ctx *Ctx) Query(key string) string {
	return ctx.r.URL.Query().Get(key)
}

func (ctx *Ctx) QueryInt(key string) int64 {
	val, err := strconv.ParseInt(ctx.r.URL.Query().Get(key), 10, 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return val
}

func (ctx *Ctx) Body() map[string]interface{} {
	bytes, _ := ioutil.ReadAll(ctx.r.Body)
	body := make(map[string]interface{})
	err := json.Unmarshal(bytes, &body)
	if err != nil {
		log.Println(err)
	}
	return body
}

func (ctx *Ctx) GetBody(receiver interface{}) interface{} {
	bytes, _ := ioutil.ReadAll(ctx.r.Body)
	err := json.Unmarshal(bytes, receiver)
	if err != nil {
		log.Println(err)
	}
	return receiver
}

type Handle func(ctx *Ctx)

type Middleware func(Handle) Handle

var routes = make(map[string]Handle)

func (s *Server) GET(uri string, handle Handle, middlewares ...Middleware) {
	s.SetRequest("GET", uri, handle, middlewares...)
}

func (s *Server) PUT(uri string, handle Handle, middlewares ...Middleware) {
	s.SetRequest("PUT", uri, handle, middlewares...)
}

func (s *Server) POST(uri string, handle Handle, middlewares ...Middleware) {
	s.SetRequest("POST", uri, handle, middlewares...)
}

func (s *Server) DELETE(uri string, handle Handle, middlewares ...Middleware) {
	s.SetRequest("DELETE", uri, handle, middlewares...)
}

func (s *Server) SetRequest(method string, uri string, handle Handle, middlewares ...Middleware) {
	handle = chain(handle, middlewares...)
	routeId := method + uri
	routes[routeId] = handle
}

func chain(h Handle, middlewares ...Middleware) Handle {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func (s *Server) Listen(port int) {
	http.HandleFunc("/", buildRouter)
	p := strconv.Itoa(port)
	fmt.Println("server is listening in http://localhost:" + p)
	http.ListenAndServe(":"+p, nil)
}

func buildRouter(w http.ResponseWriter, r *http.Request) {
	routeId := r.Method + r.URL.Path
	handle, ok := routes[routeId]
	if !ok {
		http.NotFound(w, r)
		return
	}

	ctx := Ctx{w: w, r: r}
	ctx.init()
	handle(&ctx)
}
