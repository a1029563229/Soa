# Go 版 Koa —— Soa

## Usage

- basic
```go
func main() {
  app := new(soa.Server)
  app.GET("/", hello.Hello)
}
```

- middleware
```go
func mid() soa.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("mid")
			next(w, r)
		}
	}
}

func main() {
  app := new(soa.Server)
  app.GET("/", hello.Hello, mid())
}
```

