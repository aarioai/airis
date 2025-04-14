#!/bin/bash
set -euo pipefail

. /opt/aa/lib/aa-posix-lib.sh

# cmd.sh
# 项目内快捷使用的脚本，可以将本脚本复制至 project/xxx/cmd/cmd.sh， 方便使用

PROTOC_VERSION="29.3"

AA_GITHUB="github.com/aarioai"
MOD_ROOT="${GOPATH}/pkg/mod/${AA_GITHUB}"

CUR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CUR

ROOT_DIR="$(cd "${CUR}/.." && pwd)"
readonly ROOT_DIR

usage(){
    cat << EOF
Usage: $0 <command>
    $0 createapp <app name>
    $0 protoc [libprotoc version]
Example:
    $0 createapp testapp
    $0 protoc 29.3
EOF
    exit 1
}

declare repo_cmd

findRepoCmdDir(){
    if [ ! -f "${ROOT_DIR}/go.mod" ];then
        go mod tidy
        panic "missing ${ROOT_DIR}/go.mod"
    fi
    airis="$(grep "${AA_GITHUB}/airis v" "${ROOT_DIR}/go.mod")"
    version=${airis##* v}
    [ -n "$version" ] || panic "go.mod missing ${AA_GITHUB}/airis"

    repo_cmd="${MOD_ROOT}/airis@v${version}/cmd"
    if [ ! -d "${repo_cmd}" ]; then
        go mod tidy
    fi
 }

main(){
    [ "$#" -ge 1 ] || usage
    local cmd="$1"
    shift
    local args=("$@")

    findRepoCmdDir
    case "$cmd" in
        new|createapp|createapp.sh)
            info "${repo_cmd}/project/createapp.sh $ROOT_DIR ${args[*]}"
            "${repo_cmd}"/project/createapp.sh "$ROOT_DIR" "${args[@]}"
            ;;
        protoc|protoc.sh)
            local protoc_version="${1-"$PROTOC_VERSION"}"
            info "${repo_cmd}/project/protoc.sh $ROOT_DIR $protoc_version"
            "${repo_cmd}"/project/protoc.sh "$ROOT_DIR" "$protoc_version"
            ;;
        *)
            panic "invalid command: ${cmd}"
    esac
}

main "$@"