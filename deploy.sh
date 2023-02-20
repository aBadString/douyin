# 仅在第一次初始化时执行, 创建运行时镜像和编译容器
function init() {
    # 创建带 ffmpeg 环境的镜像
    docker build -f docker/Dockerfile_runtime -t abadstring/debian_ffmpeg .

    # 创建编译镜像和容器
    docker build -f docker/Dockerfile_compile -t abadstring/golang_compile:1.19.6 .
    docker create --name douyin-compile \
           -v $PWD/src:/go/src/douyin \
           -v $PWD/bin:/go/bin/douyin \
           abadstring/golang_compile:1.19.6
    # 编译脚本
    echo -e \
        "go mod tidy\n \
        go build -o /go/bin/douyin  # /go/bin 在环境变量 PATH 中" \
        > ./src/compile.sh
}

# 每次需要发布新版本时执行, 编译源代码和创建应用镜像
function compile() {
    git pull
    # 使用编译容器编译源代码, 编译后的可执行文件在 /bin 下
    docker start douyin-compile
    docker-compose build
}

# 重新创建并运行新的应用镜像的容器
function run() {
    docker-compose up -d
}


case $1 in
  "init")
    init
    compile
    run
    ;;
  "compile")
    compile
    run
    ;;
  "clean")
    docker-compose down
    docker rm douyin-compile
    ;;
  *)
    echo "使用方法: bash deploy.sh [init|compile|clean]"
    ;;
esac