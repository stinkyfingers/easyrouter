package easyrouter

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/stinkyfingers/prefixtree"
	"golang.org/x/net/websocket"
)

type Server struct {
	// routemap     map[string]map[string]Route
	tree         map[string]*prefixtree.Node
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

	s.tree = make(map[string]*prefixtree.Node)
	for i, route := range s.Routes {
		if len(route.Methods) == 0 {
			route.Methods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "PATCH"}
		}
		for _, method := range route.Methods {
			if _, ok := s.tree[method]; !ok {
				s.tree[method] = prefixtree.NewTree()
			}
			cleanedPath := getPathAsWildcard(route.Path)
			s.tree[method].Insert(cleanedPath, i)
		}
	}
}

func getPathAsWildcard(path string) string {
	arr := strings.Split(path, "/")
	var builder strings.Builder
	for _, s := range arr {
		if s == "" {
			continue
		}
		if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			builder.WriteString("/*")
			continue
		}
		builder.WriteString("/" + s)
	}
	return builder.String()
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
	var route Route
	if methodMap, ok := s.tree[r.Method]; !ok {
		return s.DefaultRoute
	} else {
		cleanedPath := getPathAsWildcard(r.URL.Path)
		res, routeIndex := methodMap.Find(cleanedPath)
		if !res {
			return route
		}
		return s.Routes[routeIndex]
	}
	return route
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
