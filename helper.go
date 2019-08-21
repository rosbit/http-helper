package helper

import (
	"github.com/urfave/negroni"
	"github.com/gernest/alien"
	"net/http"
)

type HttpHelper struct {
	n *negroni.Negroni
	m *alien.Mux
}

func NewHelper(handlers ...negroni.Handler) *HttpHelper {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	if handlers != nil {
		for _, handler := range handlers {
			n.Use(handler)
		}
	}
	m := alien.New()
	n.UseHandler(m)
	return &HttpHelper{n, m}
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

