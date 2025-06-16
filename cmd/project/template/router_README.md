# Iris Router

## Iris Router Middleware App Methods

| App Method | 	Timing             | 	Scope	                   | Executions	        | Typical Use Cases                     |
|------------|---------------------|---------------------------|--------------------|---------------------------------------|
| UseGlobal  | 	Earliest           | 	All requests	            | Every request	     | Global logging, basic auth            |
| UseOnce	   | Registration order	 | Scoped                    | 	Once per request	 | Initialization, DB connections        |
| UseFunc	   | Registration order  | 	Scoped                   | 	Every request     | 	Third-party middleware integration   |
| Use        | 	Pre-matching	      | Current group + children	 | Every request      | 	Group auth, pre-processing           |
| UseRouter	 | Post-matching       | 	Matched routes only	     | Every request	     | Route-specific logic, post-processing |

```go
app := iris.New()

app.UseOnce(func(ctx iris.Context) {
    println("2. UseOnce (single execution)")
    ctx.Next()
})

app.Use(func(ctx iris.Context) {
    println("3. Use (pre-matching)")
    ctx.Next()
})

app.UseRouter(func(ctx iris.Context) {
    println("4. UseRouter (post-matching)")
    ctx.Next()
})
app.UseGlobal(func(ctx iris.Context) {
println("1. UseGlobal (all requests)")
ctx.Next()
})


app.Get("/", func(ctx iris.Context) {
    println("5. Handler")
    ctx.Text("Hello World")
})
```

**Request handling sequence:**
1. UseGlobal
2. UseOnce
3. Use
4. UseRouter
5. Handler


## Iris Router Middleware Party Methods

* **App Middleware**
  - app.Use() - runs before any route matching 
  - app.UseRouter() - runs after route matching
* **Party (Route Group) Middleware**
  - Party.Use() - runs for all routes in the group (before matching)
  - Party.UseRouter() - runs only for matched routes in the group

App Use → Party Use → App UseRouter → Party UseRouter → Handler

```go
app := iris.New()
api := app.Party("/api")

api := app.Party("/api")


api.Use(func(ctx iris.Context) {
    println("2. Party.Use middleware")
    ctx.Next()
})

api.UseRouter(func(ctx iris.Context) {
    println("4. Party.UseRouter middleware")
    ctx.Next()
})

api.Get("/test", func(ctx iris.Context) {
    println("5. Route handler")
    ctx.Text("OK")
})

app.Use(func(ctx iris.Context) {
    println("1. Global Use middleware")
    ctx.Next()
})

app.UseRouter(func(ctx iris.Context) {
    println("3. Global UseRouter middleware")
    ctx.Next()
})


api2 := app.Party("/api")
api2.Use(basicauth.Default(...))
api.Get("/hello", ...)
```

**Request `/api/test` execution sequences:**
1. Global Use middleware
2. Party.Use middleware
3. Global UseRouter middleware
4. Party.UseRouter middleware
5. Route handler

Note: This request won't go through basic auth because:
* It matches the first /api group's /test route
* Iris uses first-matched route and won't check subsequent groups with same prefix

**Request `/api/hello` execution sequence:**
1. Global Use middleware
2. Party.Use basic auth
3. Global UseRouter middleware
4. Route handler hello