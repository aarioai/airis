# Project Structure

**Project Overview**
* app: application
  * app_name1
  * app_name2
  * ...
* boot: bootstrap
* cmd: command
* config: configuration of ini/rsa files
* docs
* frontend: source code of Javascript/CSS assets and views
* repair
* sdk
* storage: log files and other static assets
* tests: unit tests
* go.mod
* main.go

**APP Overview**
* bo: business object
* cache
* conf: configuration of go files
* dic
* entity
  * mo: Mongodb Object
  * po: Persistent Object,   entities --> po --> model/service
  * xxx.go: SQL entity
* enum
* job
* module
  * bs: B/S APIs
    * controller
    * dto: Data Transfer Object
    * model
    * phone/pc/view: view controller
    * service.go
    * ...
  * cms: Content Management System
  * ss: S/S APIs
  * task
  * ...
* mservice
* private
* service

```
app
  app_name
    bo
    cache
    conf
    dic
    entity
      mo
      po
    enum
    job
    module
      bs
      cms
      ss
      task
      ...
    mservice
    private
    service
  ...
  router
    middleware  
  rpc
boot
cmd
config
docs
frontend
repair
sdk
storage
tests
go.mod
main.go    
```