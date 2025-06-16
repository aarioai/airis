# Iris 路由

## Iris App 路由中间件

| APP 方法     | 	顺序     | 	范围	        | 触发	     | 用例         |
|------------|---------|-------------|---------|------------|
| UseGlobal  | 	全局最早   | 	所有请求		     | 每次请求    | 全局日志、最基础认证 |
| UseOnce	   | 按注册顺序		 | 注册的范围       | 	仅执行一次	 | 初始化数据库连接   |
| UseFunc	   | 按注册顺序	  | 	注册的范围      | 	每次请求   | 	整合第三方中间件  |
| Use        | 	路由匹配前	 | 当前路由组及子路由		 | 每次请求    | 	路由组认证、预处理 |
| UseRouter	 | 路由匹配后   | 仅匹配成功的路由	   | 每次请求	   | 路由特定逻辑、后处理 |

```go
app := iris.New()

app.UseOnce(func(ctx iris.Context) {
    println("2. UseOnce (只会出现一次)")
    ctx.Next()
})

app.Use(func(ctx iris.Context) {
    println("3. Use")
    ctx.Next()
})

app.UseRouter(func(ctx iris.Context) {
    println("4. UseRouter")
    ctx.Next()
})

app.UseGlobal(func(ctx iris.Context) {
println("1. UseGlobal")
ctx.Next()
})

app.Get("/", func(ctx iris.Context) {
    println("5. Handler")
    ctx.Text("Hello")
})

app.Listen(":8080")
```

**请求输出顺序：**
1. UseGlobal
2. UseOnce (只会出现一次)
3. Use
4. UseRouter
5. Handler

## Iris Party 路由中间件

* **App 中间件**
  - app.Use() - 每个路由匹配前执行
  - app.UseRouter() - 每个路由匹配后执行
* **Party （相当于路由组）中间件**
  - Party.Use() - 当前路由组匹配前执行
  - Party.UseRouter() - 当前路由组匹配后执行

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

**请求 `/api/test` 执行顺序：**
1. Global Use middleware
2. Party.Use middleware
3. Global UseRouter middleware
4. Party.UseRouter middleware
5. Route handler

注意：`/api/test` 不会使用后面注册的 basicauth
* Iris 只会匹配第一个路由组
* Iris 注册路由组，不会合并和检测party 前缀是否相同

**请求 `/api/hello` 执行顺序：**
1. Global Use middleware
2. Party.Use basic auth
3. Global UseRouter middleware
4. Route handler hello