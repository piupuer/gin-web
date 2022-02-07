FROM golang:1.17 AS gin-web
#FROM registry.cn-shenzhen.aliyuncs.com/piupuer/golang:1.17-alpine AS gin-web

RUN echo "----------------- Gin Web building(Production) -----------------"
# set environments
# enable go modules
ENV GO111MODULE=on
# set up an agent to speed up downloading resources
ENV GOPROXY=https://goproxy.cn,direct
# set app home dir
ENV APP_HOME /app/gin-web-prod

RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

# copy go.mod / go.sum to download dependent files
COPY go.mod go.sum ./
RUN go mod tidy

# copy source files
COPY . .

# save current git version
RUN chmod +x version.sh && ./version.sh

RUN go build -o main-prod .

# mysqldump need to use alpine-glibc
FROM frolvlad/alpine-glibc:alpine-3.12
#FROM registry.cn-shenzhen.aliyuncs.com/piupuer/frolvlad-alpine-glibc:alpine-3.12

# set project run mode
ENV GIN_WEB_MODE production
ENV CFG_SYSTEM_PORT 8080
ENV APP_HOME /app/gin-web-prod

RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

COPY --from=gin-web $APP_HOME/conf ./conf/
COPY --from=gin-web $APP_HOME/main-prod .
COPY --from=gin-web $APP_HOME/gitversion .

COPY docker-conf/mysql/mysqldump /usr/bin/mysqldump

# use ali apk mirros
# change timezone to Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update \
  && apk add tzdata \
  && apk add curl \
  && apk add libstdc++ \
  && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && echo "Asia/Shanghai" > /etc/timezone
# verify that the time zone has been modified
# RUN date -R

CMD ["./main-prod"]

HEALTHCHECK --interval=5s --timeout=3s CMD curl -fs http://127.0.0.1:$CFG_SYSTEM_PORT/api/ping || exit 1;
