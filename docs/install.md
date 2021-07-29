# 安装步骤

## 后端

### 1. 安装Golang环境
至少1.14版本, [点我下载](https://golang.org/dl/), 具体安装步骤不再此处过多描述

### 2. 下载项目到本地

```bash
git clone https://github.com/piupuer/gin-web
```

### 3. 使用国内代理(加速依赖下载, 能翻也可以不配置)
#### 3.1. 只针对当前会话生效

```bash
export GOPROXY=https://goproxy.cn,direct
```

#### 3.2. 写入环境变量文件(以UNIX系统为例)

```bash
# 进入home目录
cd ~
# 追加一行
echo .bashrc >> 'export GOPROXY=https://goproxy.cn,direct'
# 当前会话生效
source .bashrc
# 重启电脑, 永久生效
# reboot
```

### 4. 修改默认配置

```bash
cd gin-web/conf

# 配置文件包含redis/mysql等默认配置, 根据你的实际情况作出修改, 否则可能导致运行报错
# 每个配置都有详细的注释说明, 相信你一定能看得懂
cat config.dev.yml
```


### 5. 运行
#### 5.1. 命令行

```bash
cd gin-web

# 项目使用go mod管理工具, 运行自动下载依赖
go run main.go
```

#### 5.2. 开发工具(以Idea为例)
##### 5.2.1. 配置项目(配置好后会自动下载依赖)

GOROOT: 
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/idea_goroot.jpeg?raw=true" width="600" alt="Idea中配置GOROOT" />
</p>

GOPROXY(能翻也可以不配置): 
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/idea_goproxy.jpeg?raw=true" width="600" alt="Idea中配置Golang国内代理" />
</p>

##### 5.2.2. 运行

找到根目录下main.go文件:
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/idea_run_main.jpeg?raw=true" width="600" alt="Idea中运行main.go" />
</p>

运行成功会出现Server is running at:
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/idea_run_success.jpeg?raw=true" width="600" alt="Idea中运行main.go成功" />
</p>


## 前端

### 1. 安装Npm环境
node自带npm, [点我下载](https://nodejs.org/en/download/), 具体安装步骤不再此处过多描述

### 2. 下载项目到本地

```bash
git clone https://github.com/piupuer/gin-web-vue
```

### 3. 安装依赖

```bash
# 使用国内代理(加速依赖下载, 能翻也可以不配置)
# npm install -g cnpm --registry=https://registry.npm.taobao.org
# cnpm install
npm install
```

### 4. 修改默认配置

```bash
# 如果没有使用nginx配置反向代理, 则需要在此处修改后端接口地址
vim gin-web-vue/.env.development
# 端口改为后端真实端口:
VUE_APP_BASE_API = 'http://127.0.0.1:10000/api/v1'
# 消息中心会使用websocket:
VUE_APP_BASE_WS = 'ws://127.0.0.1:10000/api/v1'
```

### 5. 运行

```bash
npm run serve
```
运行成功会出现App running at:
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/vue_run_success.jpeg?raw=true" width="600" alt="运行vue成功" />
</p>
