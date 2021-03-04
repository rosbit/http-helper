package helper

import (
	"github.com/urfave/negroni"
	"github.com/go-zoo/bone"
	"log"
	"os"
	"net/http"
)

type HttpHelper struct {
	RouterGroup
	n *negroni.Negroni
}

func NewHelper(handlers ...negroni.Handler) *HttpHelper {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(newLogger())
	if handlers != nil {
		for _, handler := range handlers {
			n.Use(handler)
		}
	}
	m := bone.New()
	n.UseHandler(m)
	return &HttpHelper{
		RouterGroup: RouterGroup{
			basePath: "/",
			m: m,
		},
		n: n,
	}
}

func newLogger() *negroni.Logger {
	logger := &negroni.Logger{ALogger: log.New(os.Stdout, "[http-helper] ", 0)}
	logger.SetDateFormat(negroni.LoggerDefaultDateFormat)
	logger.SetFormat("{{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Request.RequestURI}}")
	return logger
}

func (h *HttpHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.n.ServeHTTP(w, r)
}

func (h *HttpHelper) Run(addr ...string) {
	h.n.Run(addr...)
}

func (h *HttpHelper) Use(handler negroni.Handler) {
	h.n.Use(handler)
}

func (h *HttpHelper) UseFunc(handlerFunc negroni.HandlerFunc) {
	h.n.UseFunc(handlerFunc)
}

func (h *HttpHelper) UseHandler(handler http.Handler) {
	h.n.Use(negroni.Wrap(handler))
}

func (h *HttpHelper) UseHandlerFunc(handlerFunc http.HandlerFunc) {
	h.n.UseHandlerFunc(handlerFunc)
}

func (hr *HttpHelper) NotFoundHandler(h http.Handler) {
	hr.m.NotFound(h)
}
