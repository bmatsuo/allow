package allow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// This example shows how to implement variations of Map with
// application-specific handling of 405 (method not allowed) responses.
func ExampleMap_custom() {
	type mymap map[string]http.HandlerFunc
	allow := func(mm mymap) *Allow {
		a := MapFunc(mm)
		ms := a.Methods()

		// app specific logic
		js, err := json.Marshal(map[string]interface{}{
			"error":       "MethodNotAllowed",
			"description": fmt.Sprintf("method must be one of %v", ms),
			"allow":       ms,
		})
		if err != nil {
			panic("marshal error")
		}
		a.NotAllowed = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Allow", strings.Join(ms, " "))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(js)
		}
		return a
	}

	mux := http.NewServeMux()
	mux.Handle("/puppies/", allow(mymap{
		"GET":  http.NotFound,
		"POST": http.NotFound,
	}))
}

// This example shows use of map literal method routing.
func ExampleMapFunc_literal() {
	mux := http.NewServeMux()
	mux.Handle("/puppies/", MapFunc(map[string]http.HandlerFunc{
		"POST": http.NotFound,
		"GET":  http.NotFound,
	}))
}

// This example shows use of fluent method routing.
func ExampleAllow_fluent() {
	mux := http.NewServeMux()
	mux.Handle("/puppies/", New().
		Allow("POST", http.NotFoundHandler()).
		AllowFunc("GET", http.NotFound))
}
