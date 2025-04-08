#!/bin/bash
set -euo pipefail

CUR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CUR
# aarioai/airis
ROOT_DIR="$(cd "${CUR}/.." && pwd)"
readonly ROOT_DIR
readonly MOD_UPDATE_FILE="${ROOT_DIR}/.modupdate"

declare comment
needCloseVPN=0
incrTag=1

readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

_log() {
    local level=$1
    local color=$2
    local message=$3
    echo -e "$(date '+%Y-%m-%d %H:%M:%S') ${color}${level:+[$level] }${message}${NC}"
}

log() {
    _log "" "" "$1"
}

info(){
    _log "info" "${GREEN}" "$1"
}

warn(){
    _log "warn" "${YELLOW}" "$1" >&2
}

panic() {
    _log "error" "${RED}" "$1" >&2
    exit 1
}

usage() {
    cat << EOF
Usage: $0 [options] [commit message]
Options:
    -u          Upgrade go.mod
    -t          Skip tag increment
    -i          Skip go mod update
    -h          Show this help message
EOF
    exit 1
}

while getopts "utih" opt; do
    case "$opt" in
        t) incrTag=0 ;;
        h) usage ;;
        *) usage ;;
    esac
done
shift $((OPTIND-1))

if [ $# -gt 0 ]; then
    comment="$1"
fi


build() {
    cd "$ROOT_DIR/cmd" || panic "failed to cd $ROOT_DIR/cmd"
    info "building project..."
    # @TODO
    go run build.go --root="$ROOT_DIR" --js="/data/Aa/proj/go/src/project/xixi/deploy/asset_src/lib_dev/aa-js/src/f_oss_filetype_readonly.js" || panic "Build failed"
}

handleUpdateMod(){
    local latest_update=''
    local today
    today="$(date +"%Y-%m-%d")"
    if [ -s "${MOD_UPDATE_FILE}" ]; then
        latest_update=$(cat "${MOD_UPDATE_FILE}")
    fi

    if [[ "$today" = "$latest_update" ]]; then
        return 0
    fi
    info "go get -u -v ./..."
    if ! go get -u -v ./... >/dev/null 2>&1; then
        warn "update go modules failed"
    fi
    sudo chmod a+rw "$MOD_UPDATE_FILE"
    info "save update mod date to $MOD_UPDATE_FILE"
    printf '%s' "$today" > "$MOD_UPDATE_FILE"
    cat "$MOD_UPDATE_FILE"
}

pushAndUpgradeMod() {
    cd "$ROOT_DIR" || panic "failed to cd $ROOT_DIR"

    handleUpdateMod

    info "go mod tidy"
    go mod tidy || panic "failed go mod tidy"

    info "go test ./..."
    go test ./... || panic "failed go test ./... failed"

    # check there are changes or not
    if [ -z "$(git status --porcelain)" ]; then
        info "No changes to commit"
        exit 0
    fi
    info "committing changes..."
    git add -A . || panic "failed git add -A ."
    git commit -m "$comment" || panic "failed git commit -m $comment"
    git push origin main || panic "failed git push origin main"

    if [ $incrTag -eq 1 ]; then
        handle_tags
    fi
}

handle_tags() {
    info "managing tags..."
    git tag -l | xargs git tag -d
    git fetch origin --prune
    latestTag=$(git describe --tags "$(git rev-list --tags --max-count=1)" 2>/dev/null || echo "")
    
    if [ -n "$latestTag" ]; then
        tag=${latestTag%.*}
        id=${latestTag##*.}
        id=$((id+1))
        newTag="$tag.$id"
        
        info "removing old tag: $latestTag"
        git tag -d "$latestTag"
        git push origin --delete tag "$latestTag"
        
        git tag "$newTag"
        git push origin --tags
        info "new tag created: $newTag"
    fi
}


unsetVPN() {
  if [[ $1 -eq 1 ]]; then
      echo "unset VPN"
      export http_proxy=""
      export https_proxy=""
      unset http_proxy
      unset https_poxy
  fi
}

setVPN() {
  if [ -n "${http_proxy:-}" ]; then
    info "proxy ${http_proxy} ${https_proxy}"
    return
  fi

  export http_proxy=http://127.0.0.1:8118
  export https_proxy=http://127.0.0.1:8118

  local http_code
  http_code=$(curl --max-time 3 -s -w '%{http_code}\n' -o /dev/null google.com)
  if [[ $http_code =~ ^[23][0-9]{2}$ ]]; then
    needCloseVPN=1
    echo "start VPN (HTTP $http_code)"
  else
    unsetVPN 1
    echo "check VPN failed (HTTP $http_code)"
  fi
}


main() {
  setVPN
  build
  pushAndUpgradeMod
  unsetVPN "$needCloseVPN"
  info "success!"
}

main