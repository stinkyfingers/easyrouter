# easyrouter

## Init an http router && run it:
```
s := easyrouter.Server{
	Port:   "8080",
	Routes: routes,
}

s.Run()
```

## Define Routes:
```
var routes = []easyrouter.Route{
	{
		Path:        "/",
		Handler:     handleDefault,
		Middlewares: []easyrouter.Middleware{myMiddleware, myMiddleware2},
	},
	{
		Path:    "/foo",
		Handler: handleFoo,
		Methods:  []string{"POST"},
		Middlewares: []easyrouter.Middleware{myMiddleware},
	},
	{
		Path:    "/bar/{id}",
		Handler: handleBar,
		Methods:  []string{"GET","OPTIONS"},
	},
}
```

## Write Middleware
- Middleware is type func(fn http.HandlerFunc) http.HandlerFunc
- Examples:

```
func myMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("--pre-handler middleware--"))
		fn(w, r)
		w.Write([]byte("--post-handler middleware--"))

	}
}
```

```
func myMiddleware2(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("--pre-handler middleware--"))
		fn(w, r)
		w.Write([]byte("--post-handler middleware--"))

	}
}
```
	
	