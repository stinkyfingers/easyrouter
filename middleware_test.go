package easyrouter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouteMiddleware(t *testing.T) {
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
	var routes = []Route{
		{
			Path:        "/",
			Middlewares: []Middleware{middleware1, middleware2},
			Handler:     handler,
		},
		{
			Path:   "/deck",
			Method: "POST",
		},
	}
	s := Server{
		Routes: routes,
	}
	s.MakeRoutemap()
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	s.ServeHTTP(resp, req)
	b, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains((string(b)), "pre-handler") || !strings.Contains((string(b)), "post-handler") || !strings.Contains((string(b)), "middleware2") {
		t.Error("expected middleware1 and middleware2")
	}
}

// TODO

// func TestServerMiddleware(t *testing.T) {
// 	var middleware1 = func(fn http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			w.Write([]byte("--pre-handler middleware--"))
// 			fn(w, r)
// 			w.Write([]byte("--post-handler middleware--"))
// 		}
// 	}
// 	var middleware2 = func(fn http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			w.Write([]byte("--middleware2--"))
// 			fn(w, r)
// 		}
// 	}
// 	var handler = func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("done"))
// 	}
// 	var routes = []Route{
// 		{
// 			Path:    "/",
// 			Handler: handler,
// 		},
// 	}
// 	s := Server{
// 		Routes:      routes,
// 		Middlewares: []Middleware{middleware1, middleware2},
// 	}
// 	s.MakeRoutemap()
// 	resp := httptest.NewRecorder()
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	s.ServeHTTP(resp, req)
// 	b, _ := ioutil.ReadAll(resp.Body)
// 	if !strings.Contains((string(b)), "pre-handler") || !strings.Contains((string(b)), "post-handler") || !strings.Contains((string(b)), "middleware2") {
// 		t.Error("expected middleware1 and middleware2")
// 	}
// }
