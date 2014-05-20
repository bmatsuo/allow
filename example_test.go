package allow

import "net/http"

// This example shows use of Allow using a Map literal.
func ExampleAllow_literal() {
	mux := http.NewServeMux()
	mux.Handle("/puppies/", MapFunc(map[string]http.HandlerFunc{
		"POST": http.NotFound,
		"GET":  http.NotFound,
	}))
}

// This example shows use of Allow using the fluent/chaining Map API.
func ExampleAllow_fluent() {
	mux := http.NewServeMux()
	mux.Handle("/puppies/", New().
		AllowFunc("POST", http.NotFound).
		AllowFunc("GET", http.NotFound))
}
