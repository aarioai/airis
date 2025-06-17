# 入门说明

该脚本会识别 project 作为项目package根路径

1. 在你项目位置，创建 project/<你的项目名>，如  project/tutorial
2. 将 github.com/aarioai/airis/cmd/cmd 脚本复制至 project/<你的项目名> 目录下
3. 进入 project/<你的项目名> ，执行 `./cmd/cmd new <APP名>`，即可自动生成项目文件夹
4. 将 go.mod module 位置修改为项目位置即可
5. 如果需要创建更多app，直接执行 `./cmd/cmd new <APP名>` 即可


运行main.go后，访问  `http://localhost:8080/ping`
```shell
go run main.go
```