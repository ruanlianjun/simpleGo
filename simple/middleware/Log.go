package middleware

import (
	"github.com/ruanlianjun/simpleGo/simple"
	"log"
	"time"
)

func StartLog() simple.HandleFunc{
	return func(c *simple.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
