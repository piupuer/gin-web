# 流程审批(以请假作为演示)

## 效果展示

### 1. 定义审批状态机(状态机=>新增, 这里以请假为例不再单独截图)

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/config.jpeg?raw=true" width="600" alt="定义审批状态机" />
</p>

### 2. 审批被拒绝

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/refused.jpeg?raw=true" width="600" alt="审批被拒绝" />
</p>

### 3. 审批结束

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/ended.jpeg?raw=true" width="600" alt="审批结束" />
</p>

### 4. 审批取消

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/cancelled.jpeg?raw=true" width="600" alt="审批取消" />
</p>

## 代码部分

### 写代码

#### 1. 定义回调函数

由于状态机不会强关联你的数据, 因此需要有自定义回调

##### 状态转换函数

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/transition.jpeg?raw=true" width="600" alt="状态转换函数" />
</p>

##### 获取提交明细

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/get_detail.jpeg?raw=true" width="600" alt="获取提交明细" />
</p>

##### 更新提交明细

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/update_detail.jpeg?raw=true" width="600" alt="更新提交明细" />
</p>

#### 2. 请假CRUD(这里简单看看model)

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/update_detail.jpeg?raw=true" width="600" alt="请假Model" />
</p>

### 调公共接口

<p align="center">
<img src="https://github.com/piupuer/gin-web-images/blob/master/docs/fsm/fsm.jpeg?raw=true" width="600" alt="调公共接口" />
</p>

#### 1. 获取当前用户待审批记录(GET /fsm/approving/list)

#### 2. 获取某记录的审批轨迹(GET /fsm/log/track)

#### 3. 获取某记录的提交人明细(GET /fsm/submitter/detail)

#### 4. 更新某记录的提交人明细(PATCH /fsm/submitter/detail)

#### 5. 发起审批, 通过和拒绝都调用该接口(PATCH /fsm/approve)

#### 6. 取消审批(PATCH /fsm/cancel)

## 最后

结合代码实战演练快速入门, 祝你好运!