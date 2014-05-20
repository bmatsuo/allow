package allow

import "net/http"

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
