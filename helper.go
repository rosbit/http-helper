package helper

import (
	"github.com/urfave/negroni"
	"github.com/go-zoo/bone"
	"log"
	"os"
	"fmt"
	"net/http"
)

type HttpHelper struct {
	RouterGroup
	n *negroni.Negroni
}

func NewHelper(handlers ...negroni.Handler) *HttpHelper {
	n := negroni.New()
	n.Use(negroni.NewRecovery())

	hasLogger := false
	for _, handler := range handlers {
		if _, ok := handler.(*negroni.Logger); ok {
			hasLogger = true
			n.Use(handler)
		}
	}
	if !hasLogger {
		n.Use(WithLogger("http-helper"))
	}
	for _, handler := range handlers {
		if _, ok := handler.(*negroni.Logger); !ok {
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

func WithLogger(name string) negroni.Handler {
	logger := &negroni.Logger{ALogger: log.New(os.Stdout, fmt.Sprintf("[%s] ", name), 0)}
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
