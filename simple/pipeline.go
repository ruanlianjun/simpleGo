package simple

type MiddlewareHandle func(c *Context,
	next func(c *Context))

type Pipeline struct {
	send    *Context
	through []MiddlewareHandle
	current int
}

func (p *Pipeline) Send(context *Context) *Pipeline {
	p.send = context
	return p
}

func (p *Pipeline) Through(middlewares []MiddlewareHandle) *Pipeline {
	p.through = middlewares
	return p
}

func (p *Pipeline) Exec() {

	if len(p.through) > p.current {
		m := p.through[p.current]
		p.current += 1
		m(p.send, func(c *Context) {
			p.Exec()
		})
	}
}

func (p *Pipeline) Then(then func(context *Context)) {
	var m MiddlewareHandle
	m = func(C *Context, next func(c *Context)) {
		then(C)
		next(C)
	}
	p.through = append(p.through, m)
	p.Exec()
}
