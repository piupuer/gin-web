# loki配置和使用

## 配置步骤

### clone一键部署脚本

```shell
git clone https://github.com/piupuer/gin-web-docker
cd gin-web-docker
chmod +x control.sh
```

### 启动loki

```shell
./control.sh loki
```

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/run.jpeg?raw=true" width="600" alt="启动loki" />
</p>

等待初始化完成即可: 
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/ps.jpeg?raw=true" width="600" alt="查看容器运行状态" />
</p>

### 登入grafana

默认用户名/密码: admin/admin
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_login.jpeg?raw=true" width="600" alt="登入grafana" />
</p>

跳过重置密码:
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_login_skip.jpeg?raw=true" width="600" alt="跳过重置密码" />
</p>

### 添加source

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_add_source.jpeg?raw=true" width="600" alt="添加source" />
</p>

选择loki: 
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_select_loki.jpeg?raw=true" width="600" alt="选择loki" />
</p>

输入网关地址: loki-gateway:3100
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_set_loki_uri.jpeg?raw=true" width="600" alt="输入网关地址" />
</p>

保存:
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_save.jpeg?raw=true" width="600" alt="保存" />
</p>

### 查看日志

进入explore: 
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_explore.jpeg?raw=true" width="600" alt="进入explore" />
</p>

选择需要查看的容器日志labels:
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/loki/grafana_labels.jpeg?raw=true" width="600" alt="选择需要查看的容器日志labels" />
</p>
