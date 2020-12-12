<h1>Gin Web Tips</h1>

> 注意事项, 遇到的一些问题在此说明

## 数据类型
- [ ] mysql的字段命名统一使用下划线命名, binlog同步的数据为驼峰, 为了避免再转换一次, redis中直接存储为驼峰命名（你也可以自行转换） 
- [ ] 针对mysql tinyint(1)同步至redis时将其转为uint类型, 因此如果项目中需要使用该类型, 在定义model时用uint替代（如果没有使用redis可忽略）
- [ ] 目前gorm v1.20.2版本中, 请使用Model().Preload().Count()替代Table().Preload().Count(), 否则可能会出现空指针异常（等后面gorm官方修复后会同步更新）
- [ ] 接收参数建议使用c.ShouldBind方法, 避免主动抛出400错误
- [ ] 接收JSON参数如需兼容字符串和数字, 例如{"name": 1}和{"name": "1"}，请使用request.ReqUint/ReqFloat64(一般是POST/PATCH请求, 开启Content-Type:application/json, GET请求通过form-mapping解析因此不支持)
