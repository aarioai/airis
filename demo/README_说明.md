# 示例说明

project 目录下，放了3个示例项目：
* [tutorial](./tutorial/README_说明.md)  生成代码目录前的项目初始化状态
* [simple](./simple/README_说明.md) 纯HTTP Server服务的项目
* [microservice](./microservice/README_说明.md) 同时提供 gRPC 微服务和HTTP Server服务的项目

项目名称，应根据业务需求具体命名。项目内，可以设置不同的APP。APP 本质上也是一种 service，是可以对外开放的 service。

因此项目目录下 /app/ 目录可以根据业务大小、关联度，设置不同的子服务，这些自服务统一由一个进程启动，这样在各微服务简单的时候，这样启动大幅度减少启动和运维成本。

后期，如果某个微服务重要度增加并已经趋于稳定，则可以直接将该 /app/xxx 单独部署为独立进程启动的微服务。