package easyrouter

import (
	"net/http"
	"net/url"
	"testing"
)

func TestRoutemap(t *testing.T) {
	var routes = []Route{
		{
			Path: "/",
		},
		{
			Path:   "/deck",
			Method: "POST",
		},
		{
			Path:   "/de",
			Method: "POST",
		},
		{
			Path:   "/deck/{id}",
			Method: "GET",
		},
		{
			Path:   "/deck/{id}",
			Method: "DEL",
		},
		{
			Path:   "/foo/{name}/bar/{id}",
			Method: "GET",
		},
		{
			Path:   "/foo",
			Method: "POST",
		},
	}

	s := Server{
		Routes: routes,
	}
	s.MakeRoutemap()

	u1, _ := url.Parse("http://localhost:8000/deck/123a")
	u2, _ := url.Parse("http://localhost:8000/deck")
	u3, _ := url.Parse("http://localhost:8080/deck/580c0fc9dd162c0f5c443a6a")
	u4, _ := url.Parse("http://localhost:8080/")
	requests := []http.Request{
		{
			Method: "GET",
			URL:    u1,
		},
		{
			Method: "POST",
			URL:    u2,
		},
		{
			Method: "GET",
			URL:    u3,
		},
		{
			Method: "GET",
			URL:    u4,
		},
	}
	for _, r := range requests {
		route := s.FindRoute(&r)
		t.Log(route)
	}

}
