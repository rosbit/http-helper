package helper

import (
	"net/http"
)

type ContextHandlerFunc func(c *Context)

func unwrap(handlerFunc ContextHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := NewHttpContext(w, r)
		handlerFunc(c)
	}
}

func (hr *HttpHelper) GET(pattern string, h func(c *Context)) error {
	return hr.Get(pattern, unwrap(h))
}

func (hr *HttpHelper) PUT(path string, h func(c *Context)) error {
	return hr.Put(path, unwrap(h))
}

func (hr *HttpHelper) POST(path string, h func(c *Context)) error {
	return hr.Post(path, unwrap(h))
}

func (hr *HttpHelper) PATCH(path string, h func(c *Context)) error {
	return hr.Patch(path, unwrap(h))
}

func (hr *HttpHelper) OPTIONS(path string, h func(c *Context)) error {
	return hr.Options(path, unwrap(h))
}

func (hr *HttpHelper) HEAD(path string, h func(c *Context)) error {
	return hr.Head(path, unwrap(h))
}

func (hr *HttpHelper) CONNECT(path string, h func(c *Context)) error {
	return hr.Connect(path, unwrap(h))
}

func (hr *HttpHelper) TRACE(path string, h func(c *Context)) error {
	return hr.Trace(path, unwrap(h))
}

func (hr *HttpHelper) DELETE(path string, h func(c *Context)) error {
	return hr.Delete(path, unwrap(h))
}
