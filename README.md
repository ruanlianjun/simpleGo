#### Use
```go
package main

import (
	"log"
	"net/http"
	"simpleGo/simple"
	"simpleGo/simple/middleware"
)

func main() {
	r := simple.New()
	r.Use(middleware.StartLog())
	v1 := r.Group("v1")
	{
		v1.GET("/demo", func(c *simple.Context) {
			c.String(http.StatusOK, "%s", "v1 demo")
		})
		admin := v1.Group("admin")
		admin.GET("login", func(c *simple.Context) {
			c.String(http.StatusOK,
				"%s",
				"v1 v2 admin login")
		})
	}

	v2 := r.Group("v2")
	{
		v2.GET("/demo", func(c *simple.Context) {
			c.String(http.StatusOK, "%s", "v2 demo")
		})
	}

	log.Panic(r.Run(":9098"))
}

```