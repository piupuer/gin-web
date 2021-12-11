# Swagger文档

## 编辑接口
注释一般建议写在/api/v1目录的各个方法下, 如

```text
// FindLeave
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Leave
// @Description FindLeave
// @Param params query request.Leave true "params"
// @Router /leave/list [GET]
```

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/swagger/v1.jpeg?raw=true" width="600" alt="编辑接口" />
</p>

## 同步到swagger

```shell
# 安装命令行工具
go get -u github.com/swaggo/swag/cmd/swag

cd gin-web
# 同步文档到swagger
swag init -o docs/swagger --pd --parseInternal 

# 成功后/docs目录下的文件会发生变化
```

## 更多参数请参考官方文档[swag](https://github.com/swaggo/swag)
