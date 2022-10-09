FROM golang as golang

# 配置模块代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

ADD . /opt

# 进入工作目录
WORKDIR /opt

# 打包 AMD64 架构
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o collector
RUN chmod a+x /opt/collector

FROM scratch

# 暴露服务端口
EXPOSE 38080

WORKDIR /opt

# 复制打包的 Go 文件到系统用户可执行程序目录下
COPY --from=golang /opt/collector /opt

# 容器启动时运行的命令
ENTRYPOINT ["/opt/collector", "-model" ]
CMD [ "config" ]
