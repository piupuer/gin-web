<h1 align="center">Gin Web</h1>

<div align="center">
由gin + gorm + jwt + casbin组合实现的RBAC权限管理脚手架Golang版, 搭建完成即可快速、高效投入业务开发
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/piupuer/gin-web" alt="Go version"/>
<img src="https://img.shields.io/badge/Gin-1.6.2-brightgreen" alt="Gin version"/>
<img src="https://img.shields.io/badge/Gorm-1.9.12-brightgreen" alt="Gorm version"/>
<img src="https://img.shields.io/github/license/piupuer/gin-web" alt="License"/>
</p>
</div>

## 特性

- `RESTful API` 设计规范
- `Gin` 一款高效的golang web框架
- `MySQL` 数据库存储
- `Jwt` 用户认证, 登入登出一键搞定
- `Casbin` 基于角色的访问控制模型(RBAC)
- `Gorm` 数据库ORM管理框架, 可自行扩展多种数据库类型(主分支已支持gorm 2.0)
- `Validator` 请求参数校验, 版本V9
- `Log` v1.2.2升级后日志支持两种常见的高性能日志 logrus / zap (移除日志写入本地文件, 强烈建议使用docker日志或其他日志收集工具)
- `Viper` 配置管理工具, 支持多种配置文件类型
- `Embed` go 1.16文件嵌入属性, 轻松将静态文件打包到编译后的二进制应用中
- `DCron` 分布式定时任务，同一task只在某台机器上执行一次(需要配置redis)
- `GoFunk` 常用工具包, 某些方法无需重复造轮子
- `FiniteStateMachine` 有限状态机, 常用于审批流程管理(没有使用工作流, 一是go的轮子太少, 二是有限状态机基本可以涵盖常用的审批流程)
- `Uploader` 大文件分块上传/多文件、文件夹上传Vue组件[vue-uploader](https://github.com/simple-uploader/vue-uploader/)
- `MessageCenter` 消息中心(websocket长连接保证实时性, 活跃用户上线时新增消息表, 不活跃用户不管, 有效降低数据量)
- `testing` 测试标准包, 快速进行单元测试
- `Grafana Loki` 轻量日志收集工具loki, 支持分布式日志收集(需要通过docker运行[gin-web-docker](https://github.com/piupuer/gin-web-docker))
- `Minio` 轻量对象存储服务(需要通过docker运行[gin-web-docker](https://github.com/piupuer/gin-web-docker))
- `Swagger` Swagger V2接口文档
- `Captcha` 密码输错次数过多需输入验证码
- `Sign` API接口签名(防重放攻击、防数据篡改)
- `Opentelemetry` 链路追踪, 快速分析接口耗时

## 中间件

- `Rate` 访问速率限制中间件 -- 限制访问流量
- `Exception` 全局异常处理中间件 -- 使用golang recover特性, 捕获所有异常, 保存到日志, 方便追溯
- `Transaction` 全局事务处理中间件 -- 每次请求无异常自动提交, 有异常自动回滚事务, 无需每个service单独调用(GET/OPTIONS跳过)
- `AccessLog` 请求日志中间件 -- 每次请求的路由、IP自动写入日志
- `Cors 跨域中间件` -- 所有请求均可跨域访问
- `Jwt` 权限认证中间件 -- 处理登录、登出、无状态token校验
- `Casbin` 权限访问中间件 -- 基于Cabin RBAC, 对不同角色访问不同API进行校验
- `Idempotence` 接口幂等性中间件 -- 保证接口不受网络波动影响而重复点击或提交(目前针对create接口加了处理，可根据实际情况更改)

## 默认菜单

- 首页
- 系统管理
    - 菜单管理
    - 角色管理
    - 用户管理
    - 接口管理
    - 数据字典
    - 操作日志
    - 消息推送
    - 机器管理
- 状态机
    - 状态机配置
    - 我的请假条
    - 待审批列表
- 上传组件
    - 上传示例1
    - 上传示例2(主要是针对ZIP压缩包上传及解压)
- 测试页面
    - 测试用例

## 在线演示

### 目前单体架构不满足业务需要, 已转至[Go Cinch](https://go-cinch.github.io/docs/#/README), 本项目不再提供在线演示.

## 快速开始

```
git clone https://github.com/piupuer/gin-web
cd gin-web
# 强烈建议使用golang官方包管理工具go mod, 无需将代码拷贝到$GOPATH/src目录下
# 确保go环境变量都配置好, 运行main文件
go run main.go
```

> 启动成功之后, 可在浏览器中输入: [http://127.0.0.1:10000/api/ping](http://127.0.0.1:10000/api/ping), 若不能访问请检查Go环境变量或数据库配置是否正确

## [文档](https://piupuer.github.io/gin-web-slate)

## 项目结构概览

```
├── api
│   └── v1 # v1版本接口目录(类似于Java中的controller), 如果有新版本可以继续添加v2/v3
├── conf # 配置文件目录(包含测试/预发布/生产环境配置参数及casbin模型配置)
├── docker-conf # docker相关配置文件
├── initialize # 数据初始化目录
│   ├── db # 数据库初始化脚本目录, 遵循sql-migrate规范
│   └── xxx.go # 包含各种需要初始化的全局变量, 如mysql/redis
├── middleware # 中间件目录
├── models # 存储层模型定义目录
├── pkg # 公共模块目录
│   ├── cache_service # redis缓存服务目录
│   ├── global # 全局变量目录
│   ├── redis # redis查询工具目录
│   ├── request # 请求相关结构体目录
│   ├── response # 响应相关结构体目录
│   ├── service # 数据DAO服务目录
│   ├── utils # 工具包目录
│   └── wechat # 微信接口目录
├── router # 路由目录
├── tests # 本地单元测试配置目录
├── upload # 上传文件默认目录
├── Dockerfile # docker镜像构建文件(生产环境)
├── Dockerfile.stage # docker镜像构建文件(预发布环境)
├── go.mod # go依赖列表
├── go.sum # go依赖下载历史
├── main.go # 程序主入口
├── README.md # 说明文档
├── TIPS.md # 个人踩坑记录
├── TODO.md # 已完成/待完成列表
```

## 前端

- 项目地址: [gin-web-vue](https://github.com/piupuer/gin-web-vue)
- 实现方式: Typescript(为什么使用它, JS的弱类型带来的问题实在不想再吐槽, TS提高效率, 反正笔者作为一枚后端用起来很舒服~)

## [注意事项](https://github.com/piupuer/gin-web/blob/master/TIPS.md)

## [TODO](https://github.com/piupuer/gin-web/blob/master/TODO.md)

## 特别感谢

前端:
<br/>
[Element UI](https://github.com/ElemeFE/element): A Vue.js 2.0 UI Toolkit for Web.
<br/>
[vue-element-admin](https://github.com/PanJiaChen/vue-element-admin): a production-ready front-end solution for admin
interfaces.
<br/>
[vue-typescript-admin-template](https://github.com/Armour/vue-typescript-admin-template): a production-ready front-end
solution for admin interfaces based on vue, typescript and UI Toolkit element-ui.
<br/>

后端:
<br/>
[Gin](https://github.com/gin-gonic/gin): a web framework written in Go (Golang).
<br/>
[gin-jwt](https://github.com/appleboy/gin-jwt): a middleware for Gin framework.
<br/>
[casbin](https://github.com/casbin/casbin): An authorization library that supports access control models like ACL, RBAC,
ABAC in Golang.
<br/>
[Gorm](https://github.com/jinzhu/gorm): The fantastic ORM library for Golang.
<br/>
[logrus](https://github.com/sirupsen/logrus): Logrus is a structured logger for Go (golang), completely API compatible with the standard library logger.
<br/>
[zap](https://github.com/uber-go/zap): Blazing fast, structured, leveled logging in Go.
<br/>
[lumberjack](https://github.com/natefinch/lumberjack): lumberjack is a log rolling package for Go.
<br/>
[viper](https://github.com/spf13/viper): Go configuration with fangs.
<br/>
[packr](https://github.com/gobuffalo/packr): The simple and easy way to embed static files into Go binaries.
<br/>
[go-funk](https://github.com/thoas/go-funk): A modern Go utility library which provides helpers (map, find, contains,
filter, ...).
<br/>
[limiter](https://github.com/ulule/limiter): Dead simple rate limit middleware for Go.
<br/>
[validator](https://github.com/go-playground/validator): Go Struct and Field validation, including Cross Field, Cross
Struct, Map, Slice and Array diving.
<br/>
[dcron](https://github.com/libi/dcron): 分布式定时任务库.
<br/>
[fsm](https://github.com/looplab/fsm): FSM is a finite state machine for Go.
<br/>
[sql-migrate](https://github.com/rubenv/sql-migrate): SQL Schema migration tool for Go. Based on gorp and goose.
<br/>

日志搜集:
<br/>
[loki](https://github.com/grafana/loki): Loki: like Prometheus, but for logs.

<br/>

下面几个类似本项目, 学习了大神的一些代码风格:
<br/>
[gin-admin](https://github.com/LyricTian/gin-admin): RBAC scaffolding based on Gin + Gorm/Mongo + Casbin + Wire.
<br/>
[gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin): Gin-vue-admin is a full-stack (frontend and backend
separation) framework designed for management system.
<br/>
[go-admin](https://github.com/wenjianzhang/go-admin): Gin + Vue + Element UI based scaffolding for front and back
separation management system.

## 互动交流

### 与作者对话

> 该项目是利用业余时间进行开发的, 开发思路参考了很多优秀的前后端框架, 结合自己的理解和实际需求, 做了改进.
> 如果觉得项目有不懂的地方或需要改进的地方, 欢迎提issue或pr!

### QQ群：943724601

<img src="https://github.com/piupuer/gin-web-images/blob/master/contact/qq_group.jpeg?raw=true" width="256" alt="QQ群" />

> 就不贴打赏二维码了, 不然显得项目很low, 如果您非要请我喝咖啡, 私信我, 哈哈哈~

## MIT License

    Copyright (c) 2021 piupuer

