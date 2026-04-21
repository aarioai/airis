#!/bin/bash
set -euo pipefail

# https://github.com/aarioai/opt
. /opt/aa/lib/aa-posix-lib.sh

ME=$(AbsPath "${BASH_SOURCE[0]}")
ROOT="$(ParentDir "$ME")"
MOD_UPDATE_FILE="${ROOT}/._update"

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
  Info "go get -u -v ./..."
  if ! go get -u -v ./... >/dev/null 2>&1; then
    Warn "update go modules failed"
  fi

  [ -f "$MOD_UPDATE_FILE" ] || touch "$MOD_UPDATE_FILE"
  [ -w "$MOD_UPDATE_FILE" ] || sudo chmod a+rw "$MOD_UPDATE_FILE"
  Info "save update mod date to $MOD_UPDATE_FILE"
  printf '%s' "$today" > "$MOD_UPDATE_FILE"
  cat "$MOD_UPDATE_FILE"
}

pushAndUpgradeMod() {
  local comment="$1"
  Info "push and upgrade go mod"
  CdOrPanic "$ROOT"

  handleUpdateMod

  Info "go mod tidy"
  [ -f "go.mod" ] || go mod init
  go mod tidy || Panic "failed go mod tidy"

  Info "go test ./..."
  go test ./... || Panic "failed go test ./... failed"

  # check there are changes or not
  if [ -z "$(git status --porcelain)" ]; then
    Info "No changes to commit"
    exit 0
  fi
  Info "committing changes..."
  git add -A . || Panic "failed git add -A ."
  git commit -m "$comment" || Panic "failed git commit -m $comment"
  git push origin main || Panic "failed git push origin main"

  Info "increase git tag and sync to remote"
  IncrRemoteGitTag -d

  Info "success. go visit: https://github.com/aarioai/airis/tags"
}

main() {
  Usage $# -eq 1 './cmd/git-push.sh <git comment>'
  pushAndUpgradeMod "$1"
}

main "$@"