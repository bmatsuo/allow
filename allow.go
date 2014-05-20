/*
Package allow provides HTTP method routing.

This API is experimental.

There are two APIs in the allow package. The Map* functions use map literal
syntax for clean, structured routes.  The Allow type has a fluent interface and
allows for http.Handler and http.HandlerFunc types to be used more freely.

Functions in the allow package issue runtime panics.
*/
package allow

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type mhandler struct {
	m string
	h http.Handler
}

// NotAllowed returns a basic, plaintext handler for responding to request
// methods not supported by a resource.
func NotAllowed(allowed ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allowed, " "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		switch n := len(allowed); n {
		case 0:
			fmt.Fprintln(w, "no methods allowed")
		case 1:
			fmt.Fprintln(w, "request method must be", allowed[0])
		case 2:
			fmt.Fprintf(w, "request method must be %s or %s\n", allowed[0], allowed[1])
		default:
			fmt.Fprintln(w, "request method must be one of", allowed)
		}
	})
}

// Allow is a type that maps HTTP request methods to handlers.
type Allow struct {
	// NotAllowed serves requests for which there is no method handler.  It is
	// the function's responsibility to set the Allow header and respond 405
	// (method not allowed).  If the function is nil a default implementation
	// is provided.
	NotAllowed http.HandlerFunc
	n          int
	mhs        []mhandler
	mset       map[string]bool
	ms         []string
}

// New allocates and returns an Allow.
func New() *Allow {
	a := new(Allow)
	a.mset = make(map[string]bool)
	return a
}

func (a *Allow) notAllowed(w http.ResponseWriter, r *http.Request) {
	h := http.Handler(a.NotAllowed)
	if h == nil {
		h = NotAllowed(a.ms...)
	}
	h.ServeHTTP(w, r)
}

// Methods returns the HTTP methods supported by a.
func (a *Allow) Methods() []string {
	return append([]string(nil), a.ms...)
}

// ServeHTTP serves the request using the handler associated with r.Method.
func (a *Allow) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mh := range a.mhs {
		if r.Method == mh.m {
			mh.h.ServeHTTP(w, r)
			return
		}
	}
	a.notAllowed(w, r)
}

// Allow makes a serve method requests using h.  Allow will panic if called
// with the same method more than once or if h is nil.  This method is not
// threadsafe.
func (a *Allow) Allow(method string, h http.Handler) *Allow {
	if a.mset[method] {
		panic("handler already defined")
	}
	a.mset[method] = true
	a.mhs = append(a.mhs, mhandler{method, h})
	a.ms = append(a.ms, method)
	sort.Strings(a.ms) // this is probably pretty inefficient
	return a
}

// AllowFunc behaves like Allow, but function literals can be passed without
// explicit conversion to http.HandlerFunc.
func (a *Allow) AllowFunc(method string, h http.HandlerFunc) *Allow {
	return a.Allow(method, h)
}

// Map allocates and returns an Allow that serves responses using the handler
// from m corresponding to the request's method.
func Map(m map[string]http.Handler) *Allow {
	a := New()
	for m, h := range m {
		a.Allow(m, h)
	}
	return a
}

// MapFunc is similar to Map but is easier to work with if dealing primarily
// with http.HandlerFunc types.
func MapFunc(m map[string]http.HandlerFunc) *Allow {
	a := New()
	for m, h := range m {
		a.Allow(m, h)
	}
	return a
}
