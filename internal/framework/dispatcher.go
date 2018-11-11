package framework

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Dispatcher struct {
	router     *mux.Router
	authorizer func(r *http.Request) bool
}

func (d *Dispatcher) Get(path string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	return d.handleFunc("GET", path, handler)
}

func (d *Dispatcher) Post(path string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	return d.handleFunc("POST", path, handler)
}

func (d *Dispatcher) handleFunc(method string, path string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	return d.router.HandleFunc(
		path,
		func(w http.ResponseWriter, r *http.Request) {
			if d.authorizer(r) {
				handler(w, r)
				return
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="SECELF"`)
			w.WriteHeader(401)
			w.Write([]byte("unauthorized"))
			return
		},
	).Methods(method)
}

func RegisterWithAuth(router *mux.Router, authorizer func(r *http.Request) bool, registeree func(dispatcher *Dispatcher)) {
	registeree(&Dispatcher{
		router:     router,
		authorizer: authorizer,
	})
}

func AcquireHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// http://blogs.msdn.com/b/ie/archive/2008/07/02/ie8-security-part-v-comprehensive-protection.aspx
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// http://blog.mozilla.com/security/2010/09/08/x-frame-options/
		w.Header().Set("X-Frame-Options", "DENY")

		h.ServeHTTP(w, r)
	})
}
