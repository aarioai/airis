#!/bin/bash

set -euo pipefail  # 启用严格模式，遇错即停



# 定义常量
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# aarioai/airis 项目根目录
readonly ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
readonly DEFAULT_BRANCH="main"
readonly DEFAULT_COMMENT="NO_COMMENT"
# 定义颜色输出
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly NC='\033[0m' # No Color

# 初始化参数
needCloseVPN=0
comment="$DEFAULT_COMMENT"
upgrade=0
incrTag=1
noUpdate=0



# 帮助信息
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

# 日志函数
log() {
    echo -e "${GREEN}$1${NC}"
}

panic() {
    echo -e "${RED}[error] $1${NC}" >&2
    exit 1
}

# 参数解析
while getopts "utih" opt; do
    case "$opt" in
        u) upgrade=1 ;;
        t) incrTag=0 ;;
        i) noUpdate=1 ;;
        h) usage ;;
        *) usage ;;
    esac
done
shift $((OPTIND-1))

# 获取提交信息
if [ $# -gt 0 ]; then
    comment="$1"
fi


# 构建函数
build() {
    cd "$ROOT_DIR/cmd" || panic "failed to cd $ROOT_DIR/cmd"
    log "building project..."
    go run build.go --root="$ROOT_DIR" --js="/data/Aa/proj/go/src/project/xixi/deploy/asset_src/lib_dev/aa-js/src/f_oss_filetype_readonly.js" || panic "Build failed"
}

# 更新并推送代码
pushAndUpgradeMod() {
    cd "$ROOT_DIR" || panic "failed to cd $ROOT_DIR"

    go mod tidy || panic "failed go mod tidy"

    # 运行单元测试
    log "go test ./..."
    go test ./... || panic "failed go test ./... failed"

    # 更新 go.mod
    if [ $upgrade -eq 1 ]; then
        log "go mod init"
        rm -f go.mod
        go mod init || panic "failed go mod init"
    fi

    # 更新依赖
    if [ $noUpdate -eq 0 ]; then
        log "updating dependencies..."
        go get -u -v ./... || panic "failed go get -u -v ./..."
    fi
    # Git 操作
    log "committing changes..."
    git add -A . || panic "failed git add -A ."
    git commit -m "$comment" || panic "failed git commit -m $comment"
    git push origin "$DEFAULT_BRANCH" || panic "failed git push origin $DEFAULT_BRANCH"

    # 处理标签
    if [ $incrTag -eq 1 ]; then
        handle_tags
    fi
}

# 处理Git标签
handle_tags() {
    log "managing tags..."
    git fetch --tags
    latestTag=$(git describe --tags "$(git rev-list --tags --max-count=1)" 2>/dev/null || echo "")
    
    if [ -n "$latestTag" ]; then
        tag=${latestTag%.*}
        id=${latestTag##*.}
        id=$((id+1))
        newTag="$tag.$id"
        
        log "removing old tag: $latestTag"
        git tag -d "$latestTag"
        git push origin --delete tag "$latestTag"
        
        git tag "$newTag"
        git push origin --tags
        log "new tag created: $newTag"
    fi
}
# 取消VPN
unsetVPN() {
  if [[ $1 -eq 1 ]]; then
      echo "unset VPN"
      export http_proxy=""
      export https_proxy=""
      unset http_proxy
      unset https_poxy
  fi
}
# 开启VPN
setVPN() {
  if [ -n "${http_proxy:-}" ]; then
    echo "proxy ${http_proxy} ${https_proxy}"
    return
  fi
  # 设置代理
  export http_proxy=http://127.0.0.1:8118
  export https_proxy=http://127.0.0.1:8118

  # 检查代理后的网络连接
  local http_code
  http_code=$(curl --max-time 3 -s -w '%{http_code}\n' -o /dev/null google.com)
  # 检查HTTP状态码，2xx和3xx都表示连接成功
  if [[ $http_code =~ ^[23][0-9]{2}$ ]]; then
    needCloseVPN=1
    echo "start VPN (HTTP $http_code)"
  else
    unsetVPN 1
    echo "check VPN failed (HTTP $http_code)"
  fi
}

# 主流程
main() {
  setVPN
  build
  pushAndUpgradeMod
  unsetVPN "$needCloseVPN"
  log "success!"
  log "use go get -u -v ./...  or -u to upgrade all dependencies maximum 1 time per day"
}

# 执行主流程
main