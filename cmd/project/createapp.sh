#!/bin/bash
set -euo pipefail

. /opt/aa/lib/aa-posix-lib.sh

# createapp.sh
# 生成app目录的脚本

CUR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CUR

readonly GLOBAL_DIRS=(
    app/router/middleware   \
    app/grpc                \
    boot                    \
    config                  \
    frontend/asset_src      \
    frontend/view_src       \
    frontend/dst            \
    maintain/repair         \
    maintain/tests          \
    sdk                     \

)
readonly APP_GLOBAL_DIRS=(
    bo                  \
    cache               \
    conf                \
    entity              \
    entity/mo           \
    entity/po           \
    enum                \
    job/queue/consumer  \
    module              \
    mservice            \
    private             \
    service
)
readonly COMMON_MODULES=(
    bs                  \
    cms                 \
    ss                  \
    task
)
readonly MODULE_DIRS=(
    controller
    dto
    model
)

createMainGo(){
    local project_root="$1"
    local demo="${CUR}/demo/main.go.demo"
    local dst="${project_root}/main.go"
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

createMiddlewareFile(){
    local project_root="$1"
    local demo="${CUR}/demo/middleware.go.demo"
    local dst="${project_root}/app/router/middleware/middleware.go"
    [ ! -f "$dst" ] || return 0
    cp "$demo" "$dst"
}

createBaseConfFile(){
    local app_root="$1"
    local app_name="$2"
    local template="${CUR}/project_template/conf_base.go.tpl"
    local dst="${app_root}/conf/base.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_NAME}}#${app_name}#g"  "$template" > "$dst"
}

createRouterFile(){
    local project_root="$1"
    local demo="${CUR}/demo/router.go.demo"
    local dst="${project_root}/app/router/router.go"
    [ ! -f "$dst" ] || return 0
    cp "$demo" "$dst"
}
createRouterEngineFile(){
    local project_root="$1"
    local demo="${CUR}/demo/router_engine.go.demo"
    local dst="${project_root}/app/router/engine.go"
    [ ! -f "$dst" ] || return 0
    cp "$demo" "$dst"
}
createCacheFile(){
    local app_root="$1"
    local app_base="$2"
    local template="${CUR}/project_template/cache.go.tpl"
    local dst="${app_root}/cache/cache.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
}

createModuleServiceFile(){
    local app_root="$1"
    local app_base="$2"
    local module="$3"
    local module_dir="$4"
    local template="${CUR}/project_template/module_service.go.tpl"
    local dst="${module_dir}/service.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g" -e "s#{{MODULE_NAME}}#${module}#g" "$template" > "$dst"
}

createModuleModelFile(){
    local app_base="$1"
    local model_dir="$2"
    local template="${CUR}/project_template/model.go.tpl"
    local dst="${model_dir}/model.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
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
    shift 2
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
        createModuleModelFile "$app_base" "${module_dir}/model"
        createModuleControllerFile "$app_base" "$module" "${module_dir}/controller"
    done
}

createBaseServiceFile(){
    local dir="$1"
    local app_base="$2"
    local pkg=${dir##*/}
    local template="${CUR}/project_template/base_service.go.tpl"
    local dst="${dir}/service.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{PACKAGE_NAME}}#${pkg}#g" -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
}



createServiceFile(){
    local app_root="$1"
    local app_base="$2"
    local template="${CUR}/project_template/service.go.tpl"
    local dst="${app_root}/service/service.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
}

createJobInitFile(){
    local app_root="$1"
    local app_name="$2"
    local template="${CUR}/project_template/job_init.go.tpl"
    local dst="${app_root}/job/init.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
}

createJobInitMongodbFile(){
    local app_root="$1"
    local app_base="$2"
    local template="${CUR}/project_template/job_init_mongodb.go.tpl"
    local dst="${app_root}/job/init_mongodb.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
}


createCommonServiceFile(){
    local dir="$1"
    local app_base="$2"
    local pkg=${dir##*/}
    local template="${CUR}/project_template/common_service.go.tpl"
    local dst="${dir}/service.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{PACKAGE_NAME}}#${pkg}#g" -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
}

createConfigFile(){
    local project_root="$1"
    local app_name="$2"
    mkdir -p "${project_root}/config"
    local template="${CUR}/project_template/app-local.ini.tpl"
    local dst="${project_root}/config/app-local.ini"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
}

createBootFiles(){
    local project_root="$1"
    local app_name="$2"
    mkdir -p "${project_root}/boot"

    # boot/init.go
    local template="${CUR}/project_template/boot_init.go.tpl"
    local dst="${project_root}/boot/init.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"

    # boot/register.go
    template="${CUR}/project_template/boot_register.go.tpl"
    dst="${project_root}/boot/register.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"

    # boot/tests/go
    template="${CUR}/project_template/boot_tests.go.tpl"
    dst="${project_root}/boot/tests.go"
    [ ! -f "$dst" ] || return 0
    sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
}

createStorage(){
    local project_root="$1"
    mkdir -p "${project_root}/storage/log"
}

goModTidy(){
    local project_root="$1"
    local project_base="$2"
    cdOrPanic "$project_root"
    if [ ! -f "go.mod" ]; then
        go mod init "$project_base"
    fi
    go mod tidy
}


main(){
    [ $# -ge 2 ] || e_usage "$0 <project root> <app name> [<module>...]${LF}Example: $0 test_app"

    local project_root="$1"
    local app_name="$2"
    shift 2
    local project_name="${project_root##*/}"
    local project_base="project/${project_name}"

    createDirs "$project_root" "${GLOBAL_DIRS[@]}"

    local app_root="${project_root}/app/${app_name}"
    local app_base="${project_base}/app/${app_name}"
    mkdir -p "$app_root"
    createDirs "$app_root" "${APP_GLOBAL_DIRS[@]}"

    createMainGo "$project_root"
    createMiddlewareFile "$project_root"
    createBaseConfFile "$app_root" "$app_name"
    createRouterFile "$project_root"
    createRouterEngineFile "$project_root"
    createCacheFile "$app_root" "$app_base"
    createModules "$app_root" "$app_base" "$@"
    createBaseServiceFile "${app_root}/private" "$app_base"
    createServiceFile "$app_root" "$app_base"

    createJobInitFile "$app_root" "$app_name"
    createJobInitMongodbFile "$app_root" "$app_base"
    createCommonServiceFile "${app_root}/job" "$app_base"
    createBaseServiceFile "${app_root}/job/queue" "$app_base"
    createCommonServiceFile "${app_root}/job/queue/consumer" "$app_base"

    createConfigFile "$project_root" "$app_name"
    createBootFiles "$project_root" "$app_name"

    createStorage "$project_root"

    goModTidy "$project_root" "$project_base"
    info "created app ${app_name} (${app_root})"
}

main "$@"