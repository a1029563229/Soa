# Go 版 Koa —— Soa

## Usage

### Basic

```go
package main

import (
	"fmt"
	"lib/soa/soa"
)

func main() {
	app := new(soa.Server)
	app.GET("/home", func(ctx *soa.Ctx) {
		ctx.End(200, "hello world")
	})
	app.Listen(8088)
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

### Convert Struct To BSON

```go
type Category struct {
	ID          primitive.ObjectID `bson:"_id"`
	Level       int64              `bson:"level"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Post        string             `bson:"post"`
	ParentId    primitive.ObjectID `bson:"parentId"`
	CreatedTime time.Time          `bson:"createdTime"`
	UpdatedTime time.Time          `bson:"updatedTime"`
}

category := Category{
	ID:          primitive.NewObjectID(),
	Level:			 1,
	CreatedTime: time.Now(),
}
bsonM := sutils.ToBson(category)
// (*mongo.Collection).InsertOne(context, bsonM)
```