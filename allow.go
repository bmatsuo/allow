/*
Package allow provides HTTP method routing.

This API is experimental.
*/
package allow

import "net/http"

type mhandler struct {
	m string
	h http.Handler
}

// Allow is a type that maps HTTP request methods to handlers.  Allow is not
// read-write threadsafe. Do not call it's Allow* methods concurrently or after
// it has started serving HTTP requests.
type Allow struct {
	n    int
	mhs  []mhandler
	mset map[string]bool
	ms   string
}

// New allocates and returns an Allow.
func New() *Allow {
	a := new(Allow)
	a.mset = make(map[string]bool)
	return a
}

func (a *Allow) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.n == 0 {
		w.Header().Set("Allow", "")
		http.Error(w, "no methods allowed", http.StatusMethodNotAllowed)
	}
	m := r.Method
	for _, mh := range a.mhs {
		if m == mh.m {
			mh.h.ServeHTTP(w, r)
			return
		}
	}
	w.Header().Set("Allow", a.ms)
	http.Error(w, "no methods allowed", http.StatusMethodNotAllowed)
}

// Allow makes a serve method requests using h.  Allow will panic if called
// with the same method more than once or if h is nil.
func (a *Allow) Allow(method string, h http.Handler) *Allow {
	if a.mset[method] {
		panic("handler already defined")
	}
	a.mset[method] = true
	a.mhs = append(a.mhs, mhandler{method, h})
	if a.ms != "" {
		a.ms += " "
	}
	a.ms += method
	return a
}

// See Allow.
func (a *Allow) AllowFunc(method string, h http.HandlerFunc) *Allow {
	return a.Allow(method, h)
}

// Allow returns an http.Handler that serves requests using the handler from m
// corresponding to the request's method.  If methods contains a nil value
// Allow will panic.  If methods is nil the returned handler allows no HTTP
// methods.
func Map(m map[string]http.Handler) http.Handler {
	a := New()
	for m, h := range m {
		a.Allow(m, h)
	}
	return a
}

func MapFunc(m map[string]http.HandlerFunc) http.Handler {
	a := New()
	for m, h := range m {
		a.Allow(m, h)
	}
	return a
}
