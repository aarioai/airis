#!/bin/bash
set -euo pipefail

. /opt/aa/lib/aa-posix-lib.sh

PROTOC_VERSION="29.3"

AA_GITHUB="github.com/aarioai"
MOD_ROOT="${GOPATH}/pkg/mod/${AA_GITHUB}"

CUR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CUR

ROOT_DIR="$(cd "${CUR}/.." && pwd)"
readonly ROOT_DIR

usage(){
    cat << EOF
Usage: $0 <project_project_root> <command>
    $0 <project_root> new|createapp <app name>
    $0 <project_root> protoc [libprotoc_version]
Example:
    $0 ../ new test
    $0 ../ protoc 29.3
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
    [ "$#" -ge 2 ] || usage
    local project_root="$1"
    local cmd="$2"
    shift 2
    local args=("$@")

    findRepoCmdDir
    case "$cmd" in
        new|createapp|createapp.sh)
            info "${repo_cmd}/project/createapp.sh $project_root ${args[*]}"
            "${repo_cmd}"/project/createapp.sh "$project_root" "${args[@]}"
            ;;
        protoc|protoc.sh)
            local protoc_version="${1-"$PROTOC_VERSION"}"
            info "${repo_cmd}/project/protoc.sh $project_root $protoc_version"
            "${repo_cmd}"/project/protoc.sh "$project_root" "$protoc_version"
            ;;
        *)
            panic "invalid command: ${cmd}"
    esac
}

main "$@"