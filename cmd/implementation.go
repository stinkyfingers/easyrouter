package main

import (
	"net/http"

	"github.com/stinkyfingers/easyrouter"
)

func main() {
	var middleware1 = func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("--pre-handler middleware--"))
			fn(w, r)
			w.Write([]byte("--post-handler middleware--"))
		}
	}
	var middleware2 = func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("--middleware2--"))
			fn(w, r)
		}
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("done"))
	}
	var routes = []easyrouter.Route{
		{
			Path:    "/",
			Handler: handler,
		},
	}
	s := easyrouter.Server{
		Routes:      routes,
		Middlewares: []easyrouter.Middleware{middleware1, middleware2},
		Port:        "8080",
	}
	s.Run()
}
