<h1>Gin Web Tips</h1>

> 注意事项, 遇到的一些问题在此说明

## 数据类型

- [ ] mysql的字段命名统一使用下划线命名, binlog同步的数据为驼峰, 为了避免再转换一次, redis中直接存储为驼峰命名（你也可以自行转换）
- [ ] 针对mysql tinyint(1)同步至redis时将其转为uint类型, 因此如果项目中需要使用该类型, 在定义model时用uint替代（如果没有使用redis可忽略）
- [ ] 目前gorm v1.20.2版本中, 请使用Model().Preload().Count()替代Table().Preload().Count(), 否则可能会出现空指针异常（等后面gorm官方修复后会同步更新）
- [ ] 接收参数建议使用c.ShouldBind方法, 避免主动抛出400错误
- [ ] 接收JSON参数如需兼容字符串和数字, 例如{"name": 1}和{"name": "1"}，请使用request.ReqUint/ReqFloat64(一般是POST/PATCH请求, 开启Content-Type:
  application/json, GET请求通过form-mapping解析因此不支持)
- [ ] gorm的unique标签最好配合index使用, 指定具体的索引名称
- [ ] gorm标签属性冒号后不能跟空格
- [ ] 使用gorm.Updates方法进行map[string]interface{}增量更新, key为蛇形而不是驼峰, 如user_mobile
- [ ] 使用gorm.Find方法建议不调用Error, 通过列表len来处理返回值
- [ ] 使用gorm.First方法建议不调用Error, 通过id>0来处理返回值
- [ ] 建议定义更新结构体时(request.UpdateXxxRequestStruct)全部使用指针变量, 避免零值被赋值到数据库, 如request.UpdateRoleRequestStruct
- [ ] gorm中的tag, column后面的值大小写敏感
- [ ] GET请求绑定时间字符串不能用models.LocalTime而是用string, c.Bind源码中使用form tag, 还没用到自定义UnmarshalJSON就已经报错
- [ ] 前后端通信JSON可通过zlib或gzip压缩, 提高网络传输效率(redis存JSON也适用, 节省存储空间)
- [ ] copier.Copy(toValue interface{}, fromValue interface{})中如果fromValue字段是指针类型, 而toValue对应字段非指针, 可能导致无法复制
- [ ] gorm tag内部标签不需要加单引号, 如gorm:"comment:'这是数据库注释'"应该直接写为gorm:"comment:这是数据库注释"
