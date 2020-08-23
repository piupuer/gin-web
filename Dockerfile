FROM golang:1.14 AS gin-web

RUN echo "----------------- 后端Gin Web构建 -----------------"
# 环境变量
# 开启go modules
ENV GO111MODULE=on
# 使用国内代理, 避免被墙资源无法访问
ENV GOPROXY=https://goproxy.cn
# 定义应用运行目录
ENV APP_HOME /app/gin-web

RUN mkdir -p $APP_HOME

# 设置运行目录
WORKDIR $APP_HOME

# 这里的根目录以docker-compose.yml配置build.context的为准
# 拷贝宿主机go.mod / go.sum文件到当前目录
COPY ./gin-web/go.mod ./gin-web/go.sum ./
# 下载依赖文件
RUN go mod download

# 拷贝宿主机全部文件到当前目录
COPY ./gin-web .

# 通过packr2将配置文件写入二进制文件
# 构建packr2
RUN cd $GOPATH/pkg/mod/github.com/gobuffalo/packr/v2@v2.8.0/packr2 && go build && chmod +x packr2
# 回到app目录运行packr2命令
RUN cd $APP_HOME && $GOPATH/pkg/mod/github.com/gobuffalo/packr/v2@v2.8.0/packr2/packr2 build

# 构建应用
RUN go build -o main .

# alpine镜像瘦身
FROM alpine:3.12

# 定义应用运行目录
ENV APP_HOME /app/gin-web

RUN mkdir -p $APP_HOME

# 设置运行目录
WORKDIR $APP_HOME

COPY --from=gin-web $APP_HOME/main .

# 拷贝mysqldump文件(binlog刷到redis会用到)
COPY ./gin-web/docker-conf/mysql/mysqldump /usr/bin/mysqldump

# alpine中缺少动态库，创建一个软链
RUN  mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# 暴露端口
EXPOSE 10000

# 启动应用(daemon off后台运行)
CMD ["./main", "-g", "daemon off;"]

