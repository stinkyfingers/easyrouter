package easyrouter

import (
	"net/http"
	"net/url"
	"testing"
)

var testRoutes = []Route{
	{
		Path: "/",
	},
	{
		Path: "/{id}",
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

func TestParams(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/")
	u1, _ := url.Parse("http://localhost:8080/123")
	u2, _ := url.Parse("http://localhost:8080/deck")
	u3, _ := url.Parse("http://localhost:8080/de")
	u4, _ := url.Parse("http://localhost:8080/deck/456")
	u5, _ := url.Parse("http://localhost:8080/deck/567")
	u6, _ := url.Parse("http://localhost:8080/foo/123/bar/456")
	u7, _ := url.Parse("http://localhost:8080/foo")
	r := []*http.Request{
		{
			URL: u,
		},
		{
			URL: u1,
		},
		{
			URL: u2,
		},
		{
			URL: u3,
		},
		{
			URL: u4,
		},
		{
			URL: u5,
		},
		{
			URL: u6,
		},
		{
			URL: u7,
		},
	}

	// testRoutes[1].GetParams(r[1])
	// t.Log(r[1].URL.Query(), testRoutes[1].Params)

	for i, route := range testRoutes {
		route.GetParams(r[i])
		t.Log(route.Params)
	}
}

func TestRoutemap(t *testing.T) {

	s := Server{
		Routes: testRoutes,
	}
	s.MakeRoutemap()

	u1, _ := url.Parse("http://localhost:8000/deck/123a")
	u2, _ := url.Parse("http://localhost:8000/deck")
	u3, _ := url.Parse("http://localhost:8080/deck/580c0fc9dd162c0f5c443a6a")
	u4, _ := url.Parse("http://localhost:8080/")
	u5, _ := url.Parse("http://localhost:8080/456")
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
		{
			Method: "GET",
			URL:    u5,
		},
	}
	for _, r := range requests {
		route := s.FindRoute(&r)
		t.Log(route)

	}
}
