#!/bin/bash
set -euo pipefail

. /opt/aa/lib/aa-posix-lib.sh

CUR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CUR
ROOT_DIR="$(cd "${CUR}/../.." && pwd)"
readonly ROOT_DIR
PROJECT_NAME="${ROOT_DIR##*/}"
readonly PROJECT_NAME
PROJECT_BASE="project/${PROJECT_NAME}"
readonly PROJECT_BASE
readonly GLOBAL_DIRS=(
    app/router          \
    boot                \
    config              \
    deploy/asset_src    \
    deploy/view_src     \
    repair              \
    sdk                 \
    tests
)
readonly APP_GLOBAL_DIRS=(
    bo                  \
    cache               \
    conf                \
    entity              \
    entity/po           \
    enum                \
    job/queue           \
    module              \
    mservice            \
    private             \
    service
)
readonly COMMON_MODULES=(
    bs                  \
    cns                 \
    ss                  \
    task
)
readonly MODULE_DIRS=(
    controller
    dto
    model
)

createMainGo(){
    local demo="${CUR}/demo/main.go.demo"
    local dst="${ROOT_DIR}/main.go"
    [ ! -f "$dst" ] || return 0
    cp "$demo" "$dst"
}

createDirs(){
    local base="$1"
    shift
    for dir in "$@"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "${base}/${dir}"
        fi
    done
}

createBaseConfFile(){
    local app_root="$1"
    local app_name="$2"
    local template="${CUR}/project_template/conf_base.go.tpl"
    local dst="${app_root}/conf/base.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_NAME}}#${app_name}#g"  "$template" > "$dst"
}

createCacheFile(){
    local app_root="$1"
    local app_base="$2"
    local driver_base="$3"
    local template="${CUR}/project_template/cache.go.tpl"
    local dst="${app_root}/cache/cache.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g" -e "s#{{DRIVER_BASE}}#${driver_base}#g"  "$template" > "$dst"
}

createModuleServiceFile(){
    local app_root="$1"
    local app_base="$2"
    local module="$3"
    local module_dir="$4"
    local template="${CUR}/project_template/module_service.go.tpl"
    local dst="${module_dir}/service.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g"  -e "s#{{MODULE_NAME}}#${module}#g" "$template" > "$dst"
}

createModuleModelFile(){
    local app_base="$1"
    local driver_base="$2"
    local model_dir="$3"
    local template="${CUR}/project_template/model.go.tpl"
    local dst="${model_dir}/model.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g"  -e "s#{{DRIVER_BASE}}#${driver_base}#g" "$template" > "$dst"
}
createModuleControllerFile(){
    local app_base="$1"
    local module="$2"
    local ctrl_dir="$3"
    local template="${CUR}/project_template/controller.go.tpl"
    local dst="${ctrl_dir}/controller.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g"  -e "s#{{MODULE_NAME}}#${module}#g" "$template" > "$dst"
}

createModules(){
    local app_root="$1"
    local app_base="$2"
    local driver_base="$3"
    shift
    shift
    shift
    local modules=("${COMMON_MODULES[@]}")
    if [ $# -gt 0 ]; then
        modules=("${modules[@]}" "$@")
    fi
    local module_root="${app_root}/module"

    createDirs "$module_root" "${modules[@]}"

    for module in "${modules[@]}"; do
        local module_dir="${module_root}/${module}"
        createDirs "$module_dir" "${MODULE_DIRS[@]}"
        createModuleServiceFile "$app_root" "$app_base" "$module" "$module_dir"
        createModuleModelFile "$app_base" "$driver_base" "${module_dir}/model"
        createModuleControllerFile "$app_base" "$module" "${module_dir}/controller"
    done
}

createBaseServiceFile(){
    local dir="$1"
    local pkg=${dir##*/}
    local template="${CUR}/project_template/base_service.go.tpl"
    local dst="${dir}/service.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{PACKAGE_NAME}}#${pkg}#g" "$template" > "$dst"
}



createServiceFile(){
    local app_root="$1"
    local app_base="$2"
    local template="${CUR}/project_template/service.go.tpl"
    local dst="${app_root}/service/service.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
}

goModTidy(){
    cdOrPanic "$ROOT_DIR"
    if [ ! -f "go.mod" ]; then
        go mod init "project/${PROJECT_NAME}"
    fi
    go mod tidy
}
main(){
    [ $# -ge 2 ] || e_usage "$0 <app name> <driver base> [<module>...]${LF}Example: $0 test_app sdk/lib/driver"
    local app_name="$1"
    local driver_base="${PROJECT_BASE}/$2"
    shift
    shift
    createDirs "$ROOT_DIR" "${GLOBAL_DIRS[@]}"

    local app_root="${ROOT_DIR}/app/${app_name}"
    local app_base="${PROJECT_BASE}/app/${app_name}"
    mkdir -p "$app_root"
    createDirs "$app_root" "${APP_GLOBAL_DIRS[@]}"

    createMainGo
    createBaseConfFile "$app_root" "$app_name"
    createCacheFile "$app_root" "$app_base" "$driver_base"
    createModules "$app_root" "$app_base" "$driver_base" "$@"
    createBaseServiceFile "${app_root}/private"
    createServiceFile "$app_root" "$app_base"

    goModTidy
}

main "$@"