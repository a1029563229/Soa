# Soa（A simple library for HTTP server）

<img align="right" width="160px" src="http://shadows-mall.oss-cn-shenzhen.aliyuncs.com/images/blogs/other/Jietu20191213-205649@2x.png">

Soa!

Let's Go!

<br/>
<br/>
<br/>
<br/>

## Usage

### Basic

```go
package main

import (
	"fmt"
	"lib/soa/soa"
)

func main() {
	app := soa.NewServer()
	app.GET("/home", func(ctx *soa.Ctx) {
		ctx.End(200, "hello world")
	})
	app.Listen(8088)
}
```

### Query Parameter
```go
app.GET("/something", GetSomething)
func GetSomething(ctx *soa.Ctx) {
	name := ctx.Query("name")

	page := ctx.QueryInt("page")
	pageSize := ctx.QueryInt("pageSize")
}
```

### Body Parameter
```go
app.POST("/something", GetSomething)
type Body struct {
	Name     string `bson:"name"`
	Page     int64  `bson:"page"`
	PageSize int64  `bson:"pageSize"`
}
func GetSomething(ctx *soa.Ctx) {
	body := Body{}
	ctx.GetBody(&body)
	// body will be filled by [post data]
}
```

### Middleware
```go
func mid() soa.Middleware {
	return func(next soa.Handle) soa.Handle {
		return func(ctx *soa.Ctx) {
			fmt.Println(ctx.Request.Method)
			next(ctx)
		}
	}
}

app.GET("/home", func(ctx *soa.Ctx) {
	ctx.End(200, "hello world")
}, mid())
```

### Set Response Header
```go
func GetSomething(ctx *soa.Ctx) {
	ctx.SetHeaders(soa.Header{
		"X-Pagination-Page":     1,
		"X-Pagination-Pagesize": 10,
		"X-Pagination-Total":    20,
	})
	ctx.SetHeader("key", "value")
	ctx.End(200, "something")
}
```