# redis一主多从哨兵

## 配置步骤

### clone一键部署脚本

```shell
git clone https://github.com/piupuer/gin-web-docker
cd gin-web-docker
chmod +x control.sh
```

### 设置主节点/从节点IP(改成你的局域网IP地址), 下面演示的是单机器多实例

```shell
export REDIS_MASTER_IP=10.13.2.252
export LOCAL_IP=10.13.2.252
# 起始端口(会自动分配各个节点对应端口, 可以设置从其它端口开始)
# export REDIS_PORT=6379
```

### 启动节点

```shell
./control.sh sentinel 3
```

等待初始化完成即可:
<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/redis/run_sentinel.jpeg?raw=true" width="600" alt="启动3哨兵节点" />
</p>

## 验证状态

### 查看容器

```shell
docker ps | grep redis
```

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/redis/ps.jpeg?raw=true" width="600" alt="查看容器运行状态" />
</p>

### 进入主节点

```shell
docker exec -it redis-master redis-cli -a 123456

keys *

# 查看主从状态
info REPLICATION

# 出现connected_slaves:2表示配置成功
```

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/redis/master.jpeg?raw=true" width="600" alt="查看主节点状态" />
</p>

### 进入从节点

```shell
docker exec -it redis-slave2 redis-cli -p 6381 -a 123456

keys *

# 查看主从状态
info REPLICATION
```

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/redis/slave.jpeg?raw=true" width="600" alt="查看从节点状态" />
</p>
