#!/bin/bash
set -euo pipefail

. /opt/aa/lib/aa-posix-lib.sh

DEFAULT_PROTOC_VERSION="29.3"

CUR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CUR

ROOT_DIR="$(cd "${CUR}/.." && pwd)"
readonly ROOT_DIR

PROTO_ROOT="${ROOT_DIR}/app/rpc"
readonly PROTO_ROOT


installProtoc(){
    local protoc_version="$1"
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

    go get -u google.golang.org/grpc
    go get -u google.golang.org/protobuf
    # Must use github.com/golang/protobuf, but not google.golang.org/protobuf
    go get -u github.com/golang/protobuf/protoc-gen-go
    go install github.com/golang/protobuf/protoc-gen-go
    protoc --version
}

parseProto(){
    local rpc_root="$1"
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
    local protoc_version=${1:"$DEFAULT_PROTOC_VERSION"}
    installProtoc "$protoc_version"
    parseProto "$PROTO_ROOT"
}
main "$@"