FROM golang:1.14

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
# 构建应用
RUN go build -o main .

#
# 暴露端口
EXPOSE 10000

# 启动应用(daemon off后台运行)
CMD ["./main", "-g", "daemon off;"]

