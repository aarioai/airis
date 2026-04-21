# Airis

Airis 是基于 Iris 开发的 Go 框架，其主要价值在于想实现一种像PHP一样方便的 Controller 层，并且能在微服务项目初期，使用一个进程就能方便启动和管理多个gRPC微服务和HTTP server服务。

Airis aconfig 里面默认使用 Consul 进行服务注册和服务发现，并且推荐将 Redis 和 MongoDB 作为基础数据源。

如果用户不想使用 MongoDB ，可以自行生成文件夹，或忽略通过 `./cmd/cmd new <app 名称>` 生成的代码 mongodb 部分即可。

同时 Airis Driver 项目（[https://github.com/aarioai/airis-driver](https://github.com/aarioai/airis-driver)），也提供了 MongoDB、MySQL、RabbitMQ、Redis 等驱动方法，方便开发者快捷操作。

## Demo 入门示例说明

Demo 代码，请参考 demo 目录下。

project 目录下，放了3个示例项目：
* tutorial  生成代码目录前的项目初始化状态
* simple 纯HTTP Server服务的项目
* microservice 同时提供 gRPC 微服务和HTTP Server服务的项目

项目名称，应根据业务需求具体命名。项目内，可以设置不同的APP。APP 本质上也是一种 service，是可以对外开放的 service。

因此项目目录下 /app/ 目录可以根据业务大小、关联度，设置不同的子服务，这些自服务统一由一个进程启动，这样在各微服务简单的时候，这样启动大幅度减少启动和运维成本。

后期，如果某个微服务重要度增加并已经趋于稳定，则可以直接将该 /app/xxx 单独部署为独立进程启动的微服务。


## 错误码以状态码形式输出

如果需要服务端将错误码以http状态码形式输出，可以有3种方法。部分运维工具检测健康状态通过http状态码，除这些健康状态接口外，其他接口不建议后端程序员直接抛出http状态码。为了方便快速定位服务端错误位置，服务端错误应该通过http状态码200，并以 response body 的 error code/message 形式抛出。

> 注意：服务端不抛出http状态码，并不意味着客户端不用处理http状态码错误。服务端前面经过网关层、k8s等层会被拦截而带上服务端或客户端错误相关的http状态码。
>> 建议将API以http 状态码输出错误，还是以body code 方式输出错误的选择权交给客户端决定

1. 客户端在 header 加 `X-Error-As-Status=true` ，服务端不用任何操作，这样服务端报错将以http状态码形式报错
2. 客户端在 URL 参数加  `x-error-as-status=true`，服务端不用任何操作，，这样服务端报错将以http状态码形式报错
3. 服务的在 writer 指定 `w.ErrorAsStatus=true`，客户端不用任何操作，服务端报错将以http状态码形式报错