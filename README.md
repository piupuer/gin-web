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
- `Gorm` 数据库ORM管理框架, 可自行扩展多种数据库类型
- `Validator` 请求参数校验, 版本V9
- `Lumberjack` 日志切割工具, 高效分离大日志文件, 按日期保存文件
- `Viper` 配置管理工具, 支持多种配置文件类型
- `Packr` 文件打包工具, 轻松将静态文件打包到编译后的二进制应用中
- `GoFunk` 常用工具包, 某些方法无需重复造轮子
- `testing` 测试标准包, 快速进行单元测试

## 中间件

- `RateLimiter` 访问速率限制中间件 -- 限制访问流量
- `Exception` 全局异常处理中间件 -- 使用golang recover特性, 捕获所有异常, 保存到日志, 方便追溯
- `Transaction` 全局事务处理中间件 -- 每次请求无异常自动提交, 有异常自动回滚事务, 无需每个service单独调用(GET/OPTIONS跳过)
- `AccessLog` 请求日志中间件 -- 每次请求的路由、IP自动写入日志
- `Cors 跨域中间件` -- 所有请求均可跨域访问
- `JwtAuth` 权限认证中间件 -- 处理登录、登出、无状态token校验
- `CasbinMiddleware` 权限访问中间件 -- 基于Cabin RBAC, 对不同角色访问不同API进行校验


## 默认菜单

- 首页
- 系统管理
  - 菜单管理
  - 角色管理
  - 用户管理
  - 接口管理
- 测试用例

## 快速开始

```
git clone https://github.com/piupuer/gin-web
cd gin-web
# 强烈建议使用golang官方包管理工具go mod, 无需将代码拷贝到$GOPATH/src目录下
# 确保go环境变量都配置好, 运行main文件
go run main.go
```

> 启动成功之后, 可在浏览器中输入: [http://127.0.0.1:10000/api/ping](http://127.0.0.1:10000/api/ping), 若不能访问请检查Go环境变量或数据库配置是否正确


## 项目结构概览

```
├── api
│   └── v1 # v1版本接口目录(类似于Java中的controller), 如果有新版本可以继续添加v2/v3
├── conf # 配置文件目录(包含测试/预发布/生产环境配置参数及casbin模型配置)
├── initialize # 数据初始化目录
├── logs # 日志文件默认目录(运行代码是生成)
├── middleware # 中间件目录
├── models # 存储层模型定义目录
├── pkg # 公共模块目录
│   ├── global # 全局变量目录
│   ├── request # 请求相关结构体目录
│   ├── request # 响应相关结构体目录
│   ├── service # 数据DAO服务目录
│   ├── utils # 工具包目录
│   └── route # 工具包目录
├── router # 路由目录
├── tests # 本地单元测试配置目录
```

## 前端

- 项目地址: [gin-web-vue](https://github.com/piupuer/gin-web-vue)
- 实现方式: Typescript

## 特别感谢

前端: 
<br/>
[Element UI](https://github.com/ElemeFE/element): A Vue.js 2.0 UI Toolkit for Web.
<br/>
[vue-element-admin](https://github.com/PanJiaChen/vue-element-admin): a production-ready front-end solution for admin interfaces.
<br/>
[vue-typescript-admin-template](https://github.com/Armour/vue-typescript-admin-template): a production-ready front-end solution for admin interfaces based on vue, typescript and UI Toolkit element-ui.
<br/>

后端:
<br/>
[Gin](https://github.com/gin-gonic/gin): a web framework written in Go (Golang).
<br/>
[gin-jwt](https://github.com/appleboy/gin-jwt): a middleware for Gin framework.
<br/>
[casbin](https://github.com/casbin/casbin): An authorization library that supports access control models like ACL, RBAC, ABAC in Golang.
<br/>
[Gorm](https://github.com/jinzhu/gorm): The fantastic ORM library for Golang.
<br/>
[zap](https://github.com/uber-go/zap): Blazing fast, structured, leveled logging in Go.
<br/>
[lumberjack](https://github.com/natefinch/lumberjack): lumberjack is a log rolling package for Go.
<br/>
[viper](https://github.com/spf13/viper): Go configuration with fangs.
<br/>
[packr](https://github.com/gobuffalo/packr): The simple and easy way to embed static files into Go binaries.
<br/>
[go-funk](https://github.com/thoas/go-funk): A modern Go utility library which provides helpers (map, find, contains, filter, ...).
<br/>
[limiter](https://github.com/ulule/limiter): Dead simple rate limit middleware for Go.
<br/>
[validator](https://github.com/go-playground/validator): Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving.

<br/>

下面几个类似本项目, 学习了大神的一些代码风格:
<br/>
[gin-admin](https://github.com/LyricTian/gin-admin): RBAC scaffolding based on Gin + Gorm/Mongo + Casbin + Wire.
<br/>
[gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin): Gin-vue-admin is a full-stack (frontend and backend separation) framework designed for management system.
<br/>
[go-admin](https://github.com/wenjianzhang/go-admin): Gin + Vue + Element UI based scaffolding for front and back separation management system.

## 互动交流

### 与作者对话

> 该项目是利用业余时间进行开发的, 开发思路参考了很多优秀的前后端框架, 结合自己的理解和自身需求, 做了改进.
> 您可以结合实际需要, 扩展业务需要, 如果您有好的idea请与我进行沟通, 一起探讨, 相互学习, 共同进步!
> 如果此项目对您提供了帮助, 也可以请作者喝杯咖啡, 嘿嘿~

<div>
<img src="https://github.com/piupuer/gin-web/blob/contact/images/ali_pay.jpeg?raw=true" width="256" alt="支付宝打赏" />
&nbsp;
&nbsp;
&nbsp;
&nbsp;
<img src="https://github.com/piupuer/gin-web/blob/contact/images/wechat_pay.jpeg?raw=true" width="256" alt="微信打赏" />
</div>

### QQ群：943724601

<img src="https://github.com/piupuer/gin-web/blob/contact/images/qq.jpeg?raw=true" width="256" alt="QQ群" />

## 提供一对一服务

> 另外, 作者有多年后端开发经验, 对前后端技术都有少许研究. 
> 如果想进一步学习但不知道怎么下手的童鞋, 在我工作之余, 可以有偿提供技术支持.
> 包括：gin-web/golang/typescript/vue, 帮助大家快速掌握和入门！

## MIT License

    Copyright (c) 2020 piupuer

