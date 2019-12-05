# Soa（A simple library for HTTP server）

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

### Convert Struct To BSON/JSON

```go
type Category struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Level       int64              `bson:"level" json:"level"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Poster      string             `bson:"poster,omitempty" json:"poster,omitempty"`
	ParentId    string             `bson:"parentId,omitempty" json:"parentId,omitempty"`
	CreatedTime int64              `bson:"createdTime" json:"createdTime"`
}

func GetSomething(ctx *soa.Ctx) {
	category := Category{
		ID:          primitive.NewObjectID(),
		Level:			 1,
		CreatedTime: time.Now(),
	}
	bsonM := ctx.BSON(category)
	// You can send formated json data to client
	ctx.Send(bsonM)
}
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