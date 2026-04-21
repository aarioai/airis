#!/bin/bash
set -euo pipefail

# https://github.com/aarioai/opt
. /opt/aa/lib/aa-posix-lib.sh

PROTOC_VERSION="29.3"

HERE="$(AbsDir "${BASH_SOURCE[0]}")"
readonly HERE

usage(){
  cat << EOF
Usage: $0 <project_project_root> <command>
    $0 <project_root> new|createapp <app name>
    $0 <project_root> protoc [libprotoc_version]
Command:
    new|createapp|createapp.sh
    protoc|protoc.sh
    build [dst]
Example:
    $0 ../ new test
    $0 ../ protoc 29.3
EOF
  exit 1
}

build(){
  local dst="${1:-main}"

  Info "go build -o $dst main.go"
  go build -o "$dst" main.go
  # get commit ID, pad with underlines when less than 40 characters
  commit_id=$(git rev-parse HEAD 2>/dev/null)
  if [ -z "$commit_id" ]; then
    Info ".git directory missing, git version omitted"
    return
  fi
  # total 41 characters
  printf '\n%s' "$(StrPadLeft "$commit_id" 40)" >> "$dst"
}

main(){
  [ "$#" -ge 2 ] || usage
  local project_root="$1"
  local cmd="$2"
  shift 2
  local args=("$@")
  Info "$cmd -> ${args[*]}"

  case "$cmd" in
    new|createapp|createapp.sh)
      Info "${HERE}/project/createapp.sh $project_root ${args[*]}"
      "${HERE}"/project/createapp.sh "$project_root" "${args[@]}"
      ;;
    protoc|protoc.sh)
      local protoc_version="${1-"$PROTOC_VERSION"}"
      Info "${HERE}/project/protoc.sh $project_root $protoc_version"
      "${HERE}"/project/protoc.sh "$project_root" "$protoc_version"
      ;;
    build)
      build "${args[@]}"
      ;;
    *)
        Panic "invalid command: ${cmd}"
  esac
}

main "$@"