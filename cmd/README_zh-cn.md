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
把 airis/cmd/cmd 脚本复制到项目 ./cmd/ 目录下
```shell
./cmd/cmd new <service_name>  # 在app目录下创建新项目，只要开放或准备开放对用户端的service，就可以被视为app
./cmd/cmd protoc [protoc_version]  # 在项目内自动查找.proto文件，并生成对应protobuf go文件的脚本
```