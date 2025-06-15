# GRPC

https://grpc.io/docs/languages/go/quickstart/

``shell
# 1. Install protobuf. See https://protobuf.dev/installation/
apt install -y protobuf-compiler
protoc --version

# 2. Install protoc-gen-go and protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Write proto file and build
cd sdk/<app>
protoc --go_out=./pb --go_opt=paths=source_relative    \
--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
./helloworld.proto

# 4. Write helloworld.go to implement proto
# 5. Register the server in register.go
# 6. Run GRPC server in server.go
``