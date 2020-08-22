package helper

import (
	"net/http"
)

func (hr *HttpHelper) Get(pattern string, h http.HandlerFunc) {
	hr.m.Get(pattern, h)
}

func (hr *HttpHelper) Put(path string, h http.HandlerFunc) {
	hr.m.Put(path, h)
}

func (hr *HttpHelper) Post(path string, h http.HandlerFunc) {
	hr.m.Post(path, h)
}

func (hr *HttpHelper) Patch(path string, h http.HandlerFunc) {
	hr.m.Patch(path, h)
}

func (hr *HttpHelper) Head(path string, h http.HandlerFunc) {
	hr.m.Head(path, h)
}

func (hr *HttpHelper) Options(path string, h http.HandlerFunc) {
	hr.m.Options(path, h)
}

//func (hr *HttpHelper) Connect(path string, h http.HandlerFunc) {
//	hr.m.Connect(path, h)
//}

//func (hr *HttpHelper) Trace(path string, h http.HandlerFunc) {
//	hr.m.Trace(path, h)
//}

func (hr *HttpHelper) Delete(path string, h http.HandlerFunc) {
	hr.m.Delete(path, h)
}

func (hr *HttpHelper) NotFoundHandler(h http.Handler) {
	hr.m.NotFound(h)
}

