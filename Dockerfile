FROM abadstring/debian_ffmpeg

ENV GIN_MODE release

# 拷贝可执行文件
WORKDIR /go/bin
COPY ./bin/douyin /go/bin

# 挂载卷
VOLUME ["/etc/douyin", "/var/lib/douyin"]

# 运行容器时执行
ENTRYPOINT ["/go/bin/douyin", "/etc/douyin/app.json"]