package main

import (
	"fmt"
	"github.com/ruanlianjun/simpleGo/simple"
	"log"
)

func main() {
	router := simple.New()

	//r := router.Group("/test").Middleware(func(c *simple.Context, next func(c *simple.Context)) {
	//	fmt.Println(1)
	//	next(c)
	//})
	router.Middleware(func(context *simple.Context, next func(context *simple.Context)) {
		fmt.Println(1)
		next(context)
	}).Middleware(func(context *simple.Context, next func(context *simple.Context)) {
		fmt.Println(1)
		next(context)
	}).GET("/test", func(c *simple.Context) {
		c.Writer.Write([]byte("test"))
	})

	router.Middleware(func(c *simple.Context, next func(c *simple.Context)) {
		fmt.Println(222)
		next(c)
	}).GET("/demo", func(c *simple.Context) {
		c.Writer.Write([]byte("demo"))
	})

	log.Panic(router.Run(":9098"))
}

func demo(c *simple.Context) {
	fmt.Println("simple demo")
}
