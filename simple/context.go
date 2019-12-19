package simple

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Method  string
	Path    string
	Params  map[string]string

	StatusCode int

	//middleware
	handlers []HandleFunc
	index    int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:     w,
		Request:    r,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: 0,
		index:      -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) PostForm(key string) string {
	return c.Request.PostForm.Get(key)
}

func (c *Context) Query(Key string) string {
	return c.Request.URL.Query().Get(Key)
}
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.Writer.Header().Set("Content-Type", "text/json")
	writer := json.NewEncoder(c.Writer)
	if err := writer.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
