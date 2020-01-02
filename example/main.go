package main

import (
	"github.com/ruanlianjun/simpleGo/simple"
	"log"
)

func main() {

	router := simple.New()

	v1 := router.Group("/v1")

	v1.Middleware(func(c *simple.Context, next func(c *simple.Context)) {
		c.Writer.Write([]byte("v1 middleware\n"))
		next(c)
	})

	v1.GET("/test", func(c *simple.Context) {
		c.Writer.Write([]byte("v1 test"))
	})

	v2 := router.Group("/v2")

	v2.Middleware(func(c *simple.Context, next func(c *simple.Context)) {
		next(c)
		c.Writer.Write([]byte("v2 middleware\n"))
	})

	v2.GET("/test", func(c *simple.Context) {
		c.Writer.Write([]byte("v2 test\n"))
	})

	log.Panic(router.Run(":9098"))
}
