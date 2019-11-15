package main

import (
	"fmt"
	"lib/soa/soa"
)

func mid() soa.Middleware {
	return func(next soa.Handle) soa.Handle {
		return func(ctx *soa.Ctx) {
			fmt.Println(ctx.Request.Method)
			next(ctx)
		}
	}
}

func main() {
	app := new(soa.Server)
	app.GET("/home", func(ctx *soa.Ctx) {
		ctx.End(200, "hello world")
	}, mid())
	app.Listen(8088)
}
