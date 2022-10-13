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
COPY --from=golang /opt/prometheus_template.yml /etc/
COPY --from=golang /opt/config.yaml /etc/   

# 给我们要传的参数一个初始值
# ENV MODELS="all"      

# 启动命令 获取config配置
ENTRYPOINT ["/opt/collector", "-model", "config" ]
# CMD  ${MODELS}

# 启动命令 全部模块启动
# ENTRYPOINT ["/opt/collector", "-model", "all" ]

