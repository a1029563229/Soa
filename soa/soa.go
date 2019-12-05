package soa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
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

func (ctx *Ctx) Send(data interface{}) {
	ctx.End(200, ctx.JSON(data))
}

func (ctx *Ctx) End(status int, message string) {
	ctx.w.Header().Set("content-type", "application/json")
	ctx.w.WriteHeader(status)
	fmt.Fprintf(ctx.w, message)
}

func (ctx *Ctx) BSON(bsonData interface{}) bson.M {
	return ToBson(bsonData)
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
		log.Println("GetBody: ", err)
	}
	return receiver
}

type Header map[string]string

func (ctx *Ctx) SetHeader(key string, value string) {
	ctx.w.Header().Set(key, value)
}

func (ctx *Ctx) SetHeaders(headers Header) {
	for key, value := range headers {
		ctx.SetHeader(key, value)
	}
}

func (ctx *Ctx) SetPageInfo(page int64, pageSize int64, total int64) {
	ctx.SetHeaders(Header{
		"X-Pagination-Page":     strconv.FormatInt(page, 10),
		"X-Pagination-Pagesize": strconv.FormatInt(pageSize, 10),
		"X-Pagination-Total":    strconv.FormatInt(total, 10),
	})
}

type Handle func(ctx *Ctx)

type Middleware func(Handle) Handle

var routes = make(map[string]Handle)

func NewServer() *Server {
	return new(Server)
}

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
	// if r.Method == "OPTIONS" {
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "*")
	// 	w.Header().Set("Access-Control-Allow-Headers", "*")
	// 	fmt.Fprintf(w, "")
	// 	return
	// }

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
