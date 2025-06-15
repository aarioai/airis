#!/bin/bash
set -euo pipefail

. /opt/aa/lib/aa-posix-lib.sh

# createapp.sh
# 生成app目录的脚本

CUR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CUR

readonly GLOBAL_DIRS=(
    app                     \
    boot                    \
    config                  \
    doc                     \
    frontend/asset_src      \
    frontend/view_src       \
    frontend/dst            \
    maintain/repair         \
    maintain/tests          \
    proto                   \
    router/middleware       \
    sdk
)
readonly APP_GLOBAL_DIRS=(
    bo                  \
    cache               \
    conf                \
    entity              \
    entity/mo           \
    entity/po           \
    enum                \
    grpc                \
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
    controller          \
    dto                 \
    model
)

createMainGo(){
    local project_root="$1"
    local template="${CUR}/template/main.go.tpl"
    local dst="${project_root}/main.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{ROOT}}#${project_root}#g"  "$template" > "$dst"
    fi
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
    local template="${CUR}/template/conf_base.go.tpl"
    local dst="${app_root}/conf/base.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g"  "$template" > "$dst"
    fi
}

createMiddlewareFile(){
    local project_root="$1"
    local demo="${CUR}/template/router_middleware.go.tpl"
    local dst="${project_root}/router/middleware/middleware.go"
    if [ ! -f "$dst" ]; then
        cp "$demo" "$dst"
    fi
}


createRouterFile(){
    local project_root="$1"
    local app_name="$2"
    local demo="${CUR}/template/router.go.tpl"
    local dst="${project_root}/router/router.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi
}

createRouterEngineFile(){
    local project_root="$1"
    local demo="${CUR}/template/router_engine.go.tpl"
    local dst="${project_root}/router/engine.go"
    if [ ! -f "$dst" ]; then
        cp "$demo" "$dst"
    fi
}

createCacheFile(){
    local app_root="$1"
    local app_base="$2"
    local template="${CUR}/template/cache.go.tpl"
    local dst="${app_root}/cache/cache.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
    fi
}

createModuleServiceFile(){
    local app_root="$1"
    local app_base="$2"
    local module="$3"
    local module_dir="$4"
    local template="${CUR}/template/module_service.go.tpl"
    local dst="${module_dir}/service.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g" -e "s#{{MODULE_NAME}}#${module}#g" "$template" > "$dst"
    fi
}

createModuleModelFile(){
    local app_base="$1"
    local model_dir="$2"
    local template="${CUR}/template/model.go.tpl"
    local dst="${model_dir}/model.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
    fi
}
createModuleControllerFile(){
    local app_base="$1"
    local module="$2"
    local ctrl_dir="$3"
    local template="${CUR}/template/controller.go.tpl"
    local dst="${ctrl_dir}/controller.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g"  -e "s#{{MODULE_NAME}}#${module}#g" "$template" > "$dst"
    fi
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
    local template="${CUR}/template/base_service.go.tpl"
    local dst="${dir}/service.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{PACKAGE_NAME}}#${pkg}#g" -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
    fi
}



createServiceFile(){
    local app_root="$1"
    local app_base="$2"
    local template="${CUR}/template/service.go.tpl"
    local dst="${app_root}/service/service.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
    fi
}



createJobInitFile(){
    local app_root="$1"
    local app_name="$2"
    local template="${CUR}/template/job_init.go.tpl"
    local dst="${app_root}/job/init.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi
}

createJobInitMongodbFile(){
    local app_root="$1"
    local app_base="$2"
    local template="${CUR}/template/job_init_mongodb.go.tpl"
    local dst="${app_root}/job/init_mongodb.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
    fi
}


createCommonServiceFile(){
    local dir="$1"
    local app_base="$2"
    local pkg=${dir##*/}
    local template="${CUR}/template/common_service.go.tpl"
    local dst="${dir}/service.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{PACKAGE_NAME}}#${pkg}#g" -e "s#{{APP_BASE}}#${app_base}#g" "$template" > "$dst"
    fi
}

createConfigFile(){
    local project_root="$1"
    local app_name="$2"
    mkdir -p "${project_root}/config"
    local template="${CUR}/template/app-local.ini.tpl"
    local dst="${project_root}/config/app-local.ini"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi
}

createBootFiles(){
    local project_root="$1"
    local app_name="$2"
    mkdir -p "${project_root}/boot"

    # boot/boot.go
    local template="${CUR}/template/boot.go.tpl"
    local dst="${project_root}/boot/boot.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi

    # boot/init.go
    local template="${CUR}/template/boot_init.go.tpl"
    local dst="${project_root}/boot/init.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi

    # boot/register.go
    template="${CUR}/template/boot_register.go.tpl"
    dst="${project_root}/boot/register.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi

    # boot/tests/go
    template="${CUR}/template/boot_tests.go.tpl"
    dst="${project_root}/boot/tests.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi
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
createJob(){
    local app_root="$1"
    local app_base="$2"
    local app_name="$3"
    createJobInitFile "$app_root" "$app_name"
    createJobInitMongodbFile "$app_root" "$app_base"
    createCommonServiceFile "${app_root}/job" "$app_base"
    createBaseServiceFile "${app_root}/job/queue" "$app_base"
    createCommonServiceFile "${app_root}/job/queue/consumer" "$app_base"
}

createGRPCServer(){
    local project_base="$1"
    local app_root="$2"
    local app_base="$3"
    local app_name="$4"
    mkdir -p "${app_root}/grpc/server/helloworld"
    createCommonServiceFile "${app_root}/grpc/server" "$app_base"


    local template="${CUR}/template/grpc/server/server.go.tpl"
    local dst="${app_root}/grpc/server/server.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi

    local template="${CUR}/template/grpc/server/register.go.tpl"
    local dst="${app_root}/grpc/server/register.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g"           \
            -e "s#{{APP_NAME}}#${app_name}#g"           \
            -e "s#{{PROJECT_BASE}}#${project_base}#g"   \
            "$template" > "$dst"
    fi

    local template="${CUR}/template/grpc/server/helloworld/helloworld.go.tpl"
    local dst="${app_root}/grpc/server/helloworld/helloworld.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_BASE}}#${app_base}#g"           \
            -e "s#{{APP_NAME}}#${app_name}#g"           \
            -e "s#{{PROJECT_BASE}}#${project_base}#g"   \
            "$template" > "$dst"
    fi
}

createGRPCClient(){
    local app_root="$1"
    local app_base="$2"
    local app_name="$3"
    mkdir -p "${app_root}/grpc/client"
    createCommonServiceFile "${app_root}/grpc/client" "$app_base"

    local template="${CUR}/template/grpc/client/client.go.tpl"
    local dst="${app_root}/grpc/client/client.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi
}

createProto(){
    local project_root="$1"
    local app_name="$2"
    mkdir -p "${project_root}/proto/${app_name}/pb"

    cp -f "${CUR}/template/proto/helloworld.proto" "${project_root}/proto/${app_name}/helloworld.proto"
    cp -f "${CUR}/template/proto/pb/helloworld.pb.go.tpl" "${project_root}/proto/${app_name}/pb/helloworld.pb.go"
    cp -f "${CUR}/template/proto/pb/helloworld_grpc.pb.go.tpl" "${project_root}/proto/${app_name}/pb/helloworld_grpc.pb.go"
}

createGRPCSDK(){
    local project_root="$1"
    local app_name="$2"
    mkdir -p "${project_root}/sdk/${app_name}"

    cp -f "${CUR}/template/sdk/service.go.tpl" "${project_root}/sdk/${app_name}/service.go"

    local template="${CUR}/template/sdk/grpc_client.go.tpl"
    local dst="${project_root}/sdk/${app_name}/grpc_client.go"
    if [ ! -f "$dst" ]; then
        sed -e "s#{{APP_NAME}}#${app_name}#g" "$template" > "$dst"
    fi
}

createGRPC(){
    local project_root="$1"
    local project_base="$2"
    local app_root="$3"
    local app_base="$4"
    local app_name="$5"

    cp -f "${CUR}/template/grpc/README.md" "${app_root}/grpc/README.md"
    createGRPCServer "$project_base" "$app_root" "$app_base" "$app_name"
    createGRPCClient "$app_root" "$app_base" "$app_name"
    createProto "$project_root" "$app_name"
    createGRPCSDK "$project_root" "$app_name"
}



main(){
    [ $# -ge 2 ] || e_usage "$0 <project root> <app name> [<module>...]${LF}Example: $0 test_app"

    local project_root="$1"  # path
    local app_name="$2"
    shift 2
    local project_name="${project_root##*/}"
    local project_base="project/${project_name}"  # for package import

    createDirs "$project_root" "${GLOBAL_DIRS[@]}"

    local app_root="${project_root}/app/${app_name}"    # for path
    local app_base="${project_base}/app/${app_name}"    # for package import
    mkdir -p "$app_root"
    createDirs "$app_root" "${APP_GLOBAL_DIRS[@]}"

    createMainGo "$project_root"
    createMiddlewareFile "$project_root"
    createBaseConfFile "$app_root" "$app_name"
    createRouterFile "$project_root" "$app_name"
    createRouterEngineFile "$project_root"
    createCacheFile "$app_root" "$app_base"
    createModules "$app_root" "$app_base" "$@"
    createBaseServiceFile "${app_root}/private" "$app_base"
    createServiceFile "$app_root" "$app_base"
    createGRPC "$project_root" "$project_base" "$app_root" "$app_base" "$app_name"
    createJob "$app_root" "$app_base" "$app_name"
    createConfigFile "$project_root" "$app_name"
    createBootFiles "$project_root" "$app_name"

    createStorage "$project_root"

    goModTidy "$project_root" "$project_base"
    info "created app ${app_name} (${app_root})"
}

main "$@"