# 命令说明

## git-push.sh 

options:
* -u    执行 go get -u -v ./...
* -t    跳过本次版本号+1
* -i    跳过 go mod tidy
* -h    帮助

```shell
git-push.sh [options] [comment]
```


## cmd

### 生成项目APP文件目录

该脚本会识别 project 作为项目package根路径

1. 在你项目位置，创建 project/<你的项目名>，如  project/tutorial
2. 将 github.com/aarioai/airis/cmd/cmd 脚本复制至 project/<你的项目名> 目录下
3. 进入 project/<你的项目名> ，执行 `./cmd/cmd new <APP名>`，即可自动生成项目文件夹
4. 将 go.mod module 位置修改为项目位置即可
5. 如果需要创建更多app，直接执行 `./cmd/cmd new <APP名>` 即可


```shell
./cmd/cmd new <service_name>  # 在app目录下创建新项目，只要开放或准备开放对用户端的service，就可以被视为app
./cmd/cmd protoc [protoc_version]  # 在项目内自动查找.proto文件，并生成对应protobuf go文件的脚本
```