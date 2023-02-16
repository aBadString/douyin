FROM golang:1.19.6

# 设置时区
RUN echo "Asia/Shanghai" > /etc/timezone && \
    rm /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata

# 设置环境变量
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

# 拷贝源代码
WORKDIR /go/src/douyin
COPY src .

# 换源, 安装 ffmpeg
RUN sed -i "s@http://deb.debian.org@http://mirrors.aliyun.com@g" /etc/apt/sources.list && \
    apt update && \
    apt install ffmpeg

# 下载依赖和编译
RUN go mod tidy
RUN go build -o /go/bin/douyin  # /go/bin 在环境变量 PATH 中

# 运行容器时执行
ENTRYPOINT ["douyin"]