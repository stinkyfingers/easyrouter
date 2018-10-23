package easyrouter

import (
	"net/http"
	"net/url"
	"reflect"
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
		Path:    "/foo",
		Methods: []string{"POST"},
	},
	{
		Path:    "/bar",
		Methods: []string{"POST"},
	},
	{
		Path:    "/foo/{id}",
		Methods: []string{"GET"},
	},
	{
		Path:    "/bar/{id}",
		Methods: []string{"DEL"},
	},
	{
		Path:    "/foo/{name}/bar/{id}",
		Methods: []string{"GET"},
	},
	{
		Path: "/foobar",
	},
}

func TestParams(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/")
	u1, _ := url.Parse("http://localhost:8080/123")
	u2, _ := url.Parse("http://localhost:8080/foo")
	u3, _ := url.Parse("http://localhost:8080/bar")
	u4, _ := url.Parse("http://localhost:8080/foo/456")
	u5, _ := url.Parse("http://localhost:8080/bar/567")
	u6, _ := url.Parse("http://localhost:8080/foo/123/bar/456")
	u7, _ := url.Parse("http://localhost:8080/foobar")

	tests := []struct {
		req      *http.Request
		route    Route
		expected []Param
	}{
		{
			req:      &http.Request{URL: u},
			route:    testRoutes[0],
			expected: nil,
		},
		{
			req:      &http.Request{URL: u1},
			route:    testRoutes[1],
			expected: []Param{{Key: "id", Value: "123"}},
		},
		{
			req:      &http.Request{URL: u2},
			route:    testRoutes[2],
			expected: nil,
		},
		{
			req:      &http.Request{URL: u3},
			route:    testRoutes[3],
			expected: nil,
		},
		{
			req:      &http.Request{URL: u4},
			route:    testRoutes[4],
			expected: []Param{{Key: "id", Value: "456"}},
		},
		{
			req:      &http.Request{URL: u5},
			route:    testRoutes[5],
			expected: []Param{{Key: "id", Value: "567"}},
		},
		{
			req:      &http.Request{URL: u6},
			route:    testRoutes[6],
			expected: []Param{{Key: "name", Value: "123"}, {Key: "id", Value: "456"}},
		},
		{
			req:      &http.Request{URL: u7},
			route:    testRoutes[7],
			expected: nil,
		},
	}

	for _, test := range tests {
		test.route.GetParams(test.req)
		if !reflect.DeepEqual(test.route.Params, test.expected) {
			t.Errorf("expected %v got %v", test.expected, test.route.Params)
		}
	}
}

func TestRoutemap(t *testing.T) {

	s := Server{
		Routes: testRoutes,
	}
	s.MakeRoutemap()

	u1, _ := url.Parse("http://localhost:8000/foo/123")
	u2, _ := url.Parse("http://localhost:8000/foo")
	u3, _ := url.Parse("http://localhost:8080/nopath/123")
	u4, _ := url.Parse("http://localhost:8080/")
	u5, _ := url.Parse("http://localhost:8080/456")
	tests := []struct {
		req          http.Request
		expectedPath string
	}{
		{
			req: http.Request{
				Method: "GET",
				URL:    u1,
			},
			expectedPath: "/foo/{id}",
		},
		{
			req: http.Request{
				Method: "POST",
				URL:    u2,
			},
			expectedPath: "/foo",
		},
		{
			req: http.Request{
				Method: "GET",
				URL:    u3,
			},
			expectedPath: "",
		},
		{
			req: http.Request{
				Method: "GET",
				URL:    u4,
			},
			expectedPath: "/",
		},
		{
			req: http.Request{
				Method: "GET",
				URL:    u5,
			},
			expectedPath: "/{id}",
		},
	}

	for i, test := range tests {
		route := s.FindRoute(&test.req)
		if route.Path != test.expectedPath {
			t.Errorf("expected path %s, got %s on test %d", test.expectedPath, route.Path, i)
		}
	}
}

func TestGetPathAsWildcard(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			"/foo/{id}",
			"/foo/*",
		}, {
			"/foo/{id}/bar/{attribute}",
			"/foo/*/bar/*",
		}, {
			"/foo",
			"/foo",
		}, {
			"",
			"",
		}, {
			"/{id}",
			"/*",
		},
	}

	for _, test := range tests {
		path := getPathAsWildcard(test.path)
		if path != test.expected {
			t.Errorf("got %s, expected %s", path, test.expected)
		}
	}
}
