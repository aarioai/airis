#!/bin/bash
set -euo pipefail

. /opt/aa/lib/aa-posix-lib.sh

# protoc.sh
# 在项目内自动查找.proto文件，并生成对应protobuf go文件的脚本

installProtoc(){
    local project_root="$1"
    local protoc_version="$2"
    if command -v protoc >/dev/null 2>&1; then
        protoc --version
        return 0
    fi
    cdOrPanic "/usr/local/src"

    local zipfile="https://github.com/protocolbuffers/protobuf/releases/download/v${protoc_version}/protoc-${protoc_version}-linux-x86_64.zip"

    curl -f -L -O "$zipfile"
    zip_filename=$(basename "$zipfile")
    rm -rf "${zip_filename}.d"
    unzip "$zip_filename" -d "${zip_filename}.d"
    sudo mv "${zip_filename}.d/bin/"* /usr/local/bin/

    cdOrPanic "$project_root"
    go get -u google.golang.org/grpc
    go get -u google.golang.org/protobuf
    # Must use github.com/golang/protobuf, but not google.golang.org/protobuf
    go get -u github.com/golang/protobuf/protoc-gen-go
    go install github.com/golang/protobuf/protoc-gen-go
    protoc --version
}

parseProto(){
    local project_root="$1"
    local rpc_root="${project_root}/app/rpc"
    cdOrPanic "$rpc_root"
    for dir in "$rpc_root"/*; do
        [ -d "$dir" ] || continue
        dir=$(realpath "$dir")
        echo "$dir"
        cdOrPanic "$dir"
        for proto in *.proto; do
            [[ -f "$proto" ]] || continue
            local dst="${dir}/pb"
            mkdir -p "$dst"
            protoc --go_out=plugins=grpc:"$dst" "$proto"
        done
    done
}


main(){
    cmd=$(basename "$0")
    [ "$#" -ge 1 ] || e_usage "$cmd <project root> [protoc version]${LF}Example: $cmd project 29.3"
    local project_root="$1"
    local protoc_version="${2-""}"
    installProtoc "$project_root" "$protoc_version"
    parseProto "$project_root"
}
main "$@"