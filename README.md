# easyrouter

## Init an http router:
```s := easyrouter.Server{
	Port:   "8080",
	Routes: routes,
}```

## Define Routes:
```var routes = []easyrouter.Route{
	{
		Path:        "/",
		Handler:     handleDefault,
		Middlewares: []easyrouter.Middleware{myMiddleware, myMiddleware2},
	},
	{
		Path:    "/foo",
		Handler: handleFoo,
		Method:  "POST",
		Middlewares: []easyrouter.Middleware{myMiddleware},
	},
	{
		Path:    "/bar/{id}",
		Handler: handleBar,
		Method:  "GET",
	},
}```

## Write Middleware
- Middleware is type func(fn http.HandlerFunc) http.HandlerFunc
- Examples:

```func myMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// pre-call middleware
		w.Write([]byte("--pre-handler middleware--"))
		fn(w, r)
		// post-call middleware
		w.Write([]byte("--post-handler middleware--"))

	}
}```

```func myMiddleware2(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// pre-call middleware
		w.Write([]byte("--pre-handler middleware--"))
		fn(w, r)
		// post-call middleware
		w.Write([]byte("--post-handler middleware--"))

	}
}```
	
	