package easyrouter

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/websocket"
)

type Server struct {
	routemap     map[string]map[string]Route
	Port         string
	DefaultRoute Route
	Routes       []Route
	Middlewares  []Middleware
}

type Route struct {
	Path        string
	Handler     http.HandlerFunc
	Middlewares []Middleware
	Methods     []string
	Params      []Param
	WSHandler   websocket.Handler
}
type Param struct {
	Key   string
	Value string
}

type Middleware func(fn http.HandlerFunc) http.HandlerFunc

func (s *Server) Run() error {
	s.MakeRoutemap()
	if s.DefaultRoute.Handler == nil {
		s.DefaultRoute = Route{Path: "/", Handler: func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("not found")) }}
	}
	return http.ListenAndServe(":"+s.Port, s.UniversalMiddleware(s))
}

func (s *Server) MakeRoutemap() {
	s.routemap = make(map[string]map[string]Route)
	for _, route := range s.Routes {
		for i := range route.Methods {
			if route.Methods[i] == "" {
				route.Methods[i] = "ANY"
			}
			if s.routemap[route.Methods[i]] == nil {
				s.routemap[route.Methods[i]] = make(map[string]Route)
			}
			paramRegex := regexp.MustCompile(`{.*?}`)
			p := paramRegex.ReplaceAllString(route.Path, "[^/]*")
			key := "^" + p + "$"
			s.routemap[route.Methods[i]][key] = route
		}
	}
	return
}

func (s *Server) AddMiddleware(route Route) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route.Handler.ServeHTTP(w, r)
	})
	for _, middle := range route.Middlewares {
		handler = middle(handler)
	}
	return handler
}

func (s *Server) UniversalMiddleware(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
	for _, middle := range s.Middlewares {
		handler = middle(handler)
	}
	return handler
}

func (s *Server) FindRoute(r *http.Request) Route {
	if methodMap, ok := s.routemap[r.Method]; ok {
		for k, route := range methodMap {
			reg := regexp.MustCompile(k)
			if reg.MatchString(r.URL.Path) {
				return route
			}
		}
	}

	if anyMethodMap, ok := s.routemap["ANY"]; ok {
		for k, route := range anyMethodMap {
			reg := regexp.MustCompile(k)
			if reg.MatchString(r.URL.Path) {
				return route
			}
		}
	}
	return s.DefaultRoute
}

func (r *Route) GetParams(req *http.Request) error {
	pathArray := strings.Split(r.Path, "/")
	reqArray := strings.Split(req.URL.Path, "/")
	if len(pathArray) != len(reqArray) {
		return errors.New("wrong number of params")
	}

	reg := regexp.MustCompile("{.*}")
	for i, path := range pathArray {
		if reg.MatchString(path) {
			key := strings.TrimSuffix(strings.TrimPrefix(path, "{"), "}")
			p := Param{
				Key:   key,
				Value: reqArray[i],
			}
			r.Params = append(r.Params, p)
		}
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// match route || default route
	route := s.FindRoute(r)
	handler := s.AddMiddleware(route)
	// Params get
	err := route.GetParams(r)
	if err != nil {
		route = s.DefaultRoute
	}

	// handler := route.Handler // ORIGINAL skips middleware

	// Params set
	values := r.URL.Query()
	for _, param := range route.Params {
		values.Add(param.Key, param.Value)

	}
	r.URL.RawQuery = values.Encode()

	if route.WSHandler != nil {
		websocket.Handler(route.WSHandler).ServeHTTP(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}
