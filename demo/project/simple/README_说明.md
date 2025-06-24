# HTTP服务入门说明


## Linux (Debian) 下运行说明

推荐使用 Linux 系统运行，本项目脚本基本都是 Linux 脚本。

1. 下载并安装Go语言程序 https://go.dev/dl/ ，具体步骤自行搜索
2. 安装完成后，可以询问deepseek，如何配置golang环境，并修改 goproxy 为中国源（若有稳定VPN环境，则不用修改）
3. 将 airis/demo 下载到你配置 GOPATH 目录下
4. 进入 project/simple 目录
5. 查看 go.mod `github.com/aarioai/airis v<版本号>` 是不是最新版本，如果不是[最新版本](https://github.com/aarioai/airis/tags)，手动改为最新版本
6. 通过SSH命令行输入并运行 `go mod tidy` 下载依赖的包，这一步如果网络不好，自行解决外网访问问题或GO PROXY源问题
7. 通过SSH命令行输入并运行 `go get -u -v ./...` 更新依赖包，这一步如果网络不好，自行解决外网访问问题或GO PROXY源问题
8. 通过SSH命令行输入并运行 `go run main.go` 即运行程序。
9. 浏览器打开 `http://localhost/v1/ping`   如果浏览器显示 "PONG" ，则表示运行成功
10. 如果端口冲突，可以修改 config/app-local.ini 里面 `port = 80` 为你本机空闲端口即可

## Windows 下运行说明

Windows 环境仅作参考，很多脚本尚未兼容 Windows

1. 下载并安装Go语言程序 https://go.dev/dl/  ，选择 Microsoft Windows 对应 .msi 版本，下载后按提示安装
2. 安装完成后，可以询问deepseek，如何配置golang环境，并修改 goproxy 为中国源（若有稳定VPN环境，则不用修改） 
3. 将 airis/demo 下载到你配置 GOPATH 目录下
4. 进入 project/simple 目录
5. 查看 go.mod `github.com/aarioai/airis v<版本号>` 是不是最新版本，如果不是[最新版本](https://github.com/aarioai/airis/tags)，手动改为最新版本
6. 通过“终端管理员”，输入并运行 `go mod tidy` 下载依赖的包，这一步如果网络不好，自行解决外网访问问题或GO PROXY源问题
7. 通过“终端管理员”，输入并运行 `go get -u -v ./...` 更新依赖包，这一步如果网络不好，自行解决外网访问问题或GO PROXY源问题
8. 通过“终端管理员”，运行 `go run main.go` 即运行程序。
9. 浏览器打开 `http://localhost/v1/ping`   如果浏览器显示 "PONG" ，则表示运行成功
10. 如果端口冲突，可以修改 config/app-local.ini 里面 `port = 80` 为你本机空闲端口即可

## 配置golang环境，并修改为中国源
安装完成后，可以询问deepseek，如何配置golang环境，并修改 goproxy 为中国源。

**1. 设置 GOPATH**

GOPATH 是你的工作目录，建议设置在用户目录下： 右键点击"此电脑" → 属性 → 高级系统设置 → 环境变量

在"用户变量"或"系统变量"中： 新建变量 GOPATH，值设为你的工作目录，例如 C:\Users\你的用户名\go ，编辑 Path 变量，添加 %GOPATH%\bin

**2. 验证 GOPATH**

打开“终端管理员”，输入命令，输出结果若为你设置的 GOPATH 工作目录，则表示设置成功
```bash
go env GOPATH
```

3. 设置 GOPROXY

在“终端管理员”输入命令：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## 生成其他APP

该 Linux 脚本会识别 project 作为项目package根路径

1. 在你项目位置，创建 project/<你的项目名>，如  project/tutorial
2. 将 github.com/aarioai/airis/cmd/cmd 脚本复制至 project/<你的项目名> 目录下
3. 进入 project/<你的项目名> ，执行 `./cmd/cmd new <APP名>`，即可自动生成项目文件夹
4. 将 go.mod module 位置修改为项目位置即可
5. 如果需要创建更多app，直接执行 `./cmd/cmd new <APP名>` 即可


运行main.go后，访问  `http://localhost/ping`
```shell
go run main.go
```