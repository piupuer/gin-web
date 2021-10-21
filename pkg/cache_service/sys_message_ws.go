package cache_service

import (
	"context"
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/ch"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	"github.com/piupuer/go-helper/pkg/middleware"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// 消息中心websocket

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// 心跳间隔
	heartBeatPeriod = 10 * time.Second

	// 最后一次活跃上线通知时间间隔
	lastActiveRegisterPeriod = 10 * time.Minute

	// 心跳最大重试次数
	HeartBeatMaxRetryCount = 3

	// 消息请求类型
	// 第1个数字: 1请求, 2响应
	// 第2个数字: 消息种类(请求和响应的消息种类没有直接联系)
	// 第3个数字: 种类排序
	// 心跳消息
	MessageReqHeartBeat string = "1-1-1"
	// 推送新消息
	MessageReqPush string = "1-2-1"
	// 批量已读
	MessageReqBatchRead string = "1-2-2"
	// 批量删除
	MessageReqBatchDeleted string = "1-2-3"
	// 全部已读
	MessageReqAllRead string = "1-2-4"
	// 全部删除
	MessageReqAllDeleted string = "1-2-5"

	// 消息响应类型(首字符为2)
	// 心跳消息
	MessageRespHeartBeat string = "2-1-1"
	// 普通消息
	MessageRespNormal string = "2-2-1"
	// 未读数
	MessageRespUnRead string = "2-3-1"
	// 用户上线
	MessageRespOnline string = "2-4-1"
)

var hub MessageHub

// 消息仓库, 用于维护整个消息中心连接
type MessageHub struct {
	lock sync.RWMutex
	// redis连接
	Service RedisService
	// 客户端用户id集合
	UserIds []uint
	// 用户最后活跃时间
	UserLastActive map[uint]int64
	// 客户端集合(用户id为每个socket key)
	Clients map[string]*MessageClient
	// 广播通道
	Broadcast *ch.Ch
	// 刷新用户消息通道
	RefreshUserMessage *ch.Ch
	// 幂等性token校验方法
	IdempotenceOps *middleware.IdempotenceOptions
}

// 消息客户端
type MessageClient struct {
	ctx context.Context
	// 当前socket key
	Key string
	// 当前socket连接实例
	Conn *websocket.Conn
	// 当前登录用户
	User models.SysUser
	// 当前登录用户ip地址
	Ip string
	// 发送消息通道
	Send *ch.Ch
	// 上次活跃时间
	LastActiveTime carbon.Carbon
	// 重试次数
	RetryCount uint
}

// 消息广播
type MessageBroadcast struct {
	response.MessageWsResp
	UserIds []uint `json:"-"`
}

// 启动消息中心仓库
func (s RedisService) StartMessageHub(idempotenceOps *middleware.IdempotenceOptions) MessageHub {
	// 初始化参数
	hub.Service = s
	hub.Clients = make(map[string]*MessageClient)
	hub.UserLastActive = make(map[uint]int64)
	hub.Broadcast = ch.NewCh()
	hub.RefreshUserMessage = ch.NewCh()
	hub.IdempotenceOps = idempotenceOps
	go hub.run()
	go hub.count()
	return hub
}

// 启动消息连接
func (s RedisService) MessageWs(ctx *gin.Context, conn *websocket.Conn, key string, user models.SysUser, ip string) {
	// 注册到消息仓库
	client := &MessageClient{
		ctx:  ctx,
		Key:  key,
		Conn: conn,
		User: user,
		Ip:   ip,
		Send: ch.NewCh(),
	}

	go client.register()
	// 监听数据的接收/发送/心跳
	go client.receive()
	go client.send()
	// go client.heartBeat()
}

// 运行仓库
func (h MessageHub) run() {
	for {
		select {
		// 广播(全部用户均可接收)
		case data := <-h.Broadcast.C:
			broadcast := data.(MessageBroadcast)
			for _, client := range h.getClients() {
				// 通知指定用户
				if utils.ContainsUint(broadcast.UserIds, client.User.Id) {
					client.Send.SafeSend(broadcast)
				}
			}
		// 刷新客户端消息
		case data := <-h.RefreshUserMessage.C:
			userIds := data.([]uint)
			// 同步用户消息
			hub.Service.mysql.SyncMessageByUserIds(userIds)
			for _, client := range h.getClients() {
				for _, id := range userIds {
					if client.User.Id == id {
						// 获取未读消息条数
						total, _ := hub.Service.mysql.GetUnReadMessageCount(id)
						// 将当前消息条数发送给用户
						msg := response.MessageWsResp{
							Type: MessageRespUnRead,
							Detail: resp.GetSuccessWithData(map[string]int64{
								"unReadCount": total,
							}),
						}
						client.Send.SafeSend(msg)
					}
				}
			}
		}
	}
}

// 活跃连接检查
func (h MessageHub) count() {
	// 创建定时器, 超出指定时间间隔
	ticker := time.NewTicker(heartBeatPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		// 到心跳检测时间
		case <-ticker.C:
			infos := make([]string, 0)
			for _, client := range h.getClients() {
				infos = append(infos, fmt.Sprintf("%d-%s", client.User.Id, client.Ip))
			}
			global.Log.Debug(h.Service.Q.Ctx, "[消息中心]当前活跃连接: %v", strings.Join(infos, ","))
		}
	}
}

// 获取client列表
func (h MessageHub) getClients() map[string]*MessageClient {
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	return hub.Clients
}

// 接收数据
func (c *MessageClient) receive() {
	defer func() {
		c.close()
		if err := recover(); err != nil {
			global.Log.Error(c.ctx, "[消息中心][接收端][%s]连接可能已断开: %v, %s", c.Key, err, string(debug.Stack()))
		}
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()

		// 记录活跃时间
		c.LastActiveTime = carbon.Now()
		c.RetryCount = 0

		if err != nil {
			panic(err)
		}
		// 解压数据
		// data := utils.DeCompressStrByZlib(string(msg))
		data := string(msg)
		global.Log.Debug(c.ctx, "[消息中心][接收端][%s]接收数据成功: %d, %s", c.Key, c.User.Id, data)
		// 数据转为json
		var r request.MessageWsReq
		utils.Json2Struct(data, &r)
		switch r.Type {
		case MessageReqHeartBeat:
			if _, ok := r.Data.(float64); ok {
				// 发送心跳
				c.Send.SafeSend(response.MessageWsResp{
					Type:   MessageRespHeartBeat,
					Detail: resp.GetSuccess(),
				})
			}
		case MessageReqPush:
			var data request.PushMessageReq
			utils.Struct2StructByJson(r.Data, &data)
			// 参数校验
			err = req.ValidateReturnErr(c.ctx, data, data.FieldTrans())
			detail := resp.GetSuccess()
			if err == nil {
				if !middleware.CheckIdempotenceToken(c.ctx, data.IdempotenceToken, *hub.IdempotenceOps) {
					err = errors.New(resp.IdempotenceTokenInvalidMsg)
				} else {
					data.FromUserId = c.User.Id
					err = hub.Service.mysql.CreateMessage(&data)
				}
			}
			if err != nil {
				detail = resp.GetFailWithMsg(err.Error())
			} else {
				// 刷新条数
				hub.RefreshUserMessage.SafeSend(hub.UserIds)
			}
			// 发送响应
			c.Send.SafeSend(response.MessageWsResp{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqBatchRead:
			var data req.Ids
			utils.Struct2StructByJson(r.Data, &data)
			err = hub.Service.mysql.BatchUpdateMessageRead(data.Uints())
			detail := resp.GetSuccess()
			if err != nil {
				detail = resp.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
			// 发送响应
			c.Send.SafeSend(response.MessageWsResp{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqBatchDeleted:
			var data req.Ids
			utils.Struct2StructByJson(r.Data, &data)
			err = hub.Service.mysql.BatchUpdateMessageDeleted(data.Uints())
			detail := resp.GetSuccess()
			if err != nil {
				detail = resp.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
			// 发送响应
			c.Send.SafeSend(response.MessageWsResp{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqAllRead:
			err = hub.Service.mysql.UpdateAllMessageRead(c.User.Id)
			detail := resp.GetSuccess()
			if err != nil {
				detail = resp.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
			// 发送响应
			c.Send.SafeSend(response.MessageWsResp{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		case MessageReqAllDeleted:
			err = hub.Service.mysql.UpdateAllMessageDeleted(c.User.Id)
			detail := resp.GetSuccess()
			if err != nil {
				detail = resp.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
			// 发送响应
			c.Send.SafeSend(response.MessageWsResp{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		}
	}
}

// 发送数据
func (c *MessageClient) send() {
	// 创建定时器, 超出指定时间间隔, 向前端发送ping消息心跳
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
		if err := recover(); err != nil {
			global.Log.Error(c.ctx, "[消息中心][发送端][%s]连接可能已断开: %v, %s", c.Key, err, string(debug.Stack()))
		}
	}()
	for {
		select {
		// 发送通道
		case msg, ok := <-c.Send.C:
			// 设定回写超时时间
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// send通道已关闭
				c.writeMessage(websocket.CloseMessage, "closed")
				panic("connection closed")
			}

			// 发送文本消息
			if err := c.writeMessage(websocket.TextMessage, utils.Struct2Json(msg)); err != nil {
				panic(err)
			}
		// 长时间无新消息
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 发送ping消息
			if err := c.writeMessage(websocket.PingMessage, "ping"); err != nil {
				panic(err)
			}
		}
	}
}

// 回写消息
func (c MessageClient) writeMessage(messageType int, data string) error {
	// 字符串压缩
	// s, _ := utils.CompressStrByZlib(data)
	s := &data
	global.Log.Debug(c.ctx, "[消息中心][发送端][%s] %v", c.Key, *s)
	return c.Conn.WriteMessage(messageType, []byte(*s))
}

// 心跳检测
func (c *MessageClient) heartBeat() {
	// 创建定时器, 超出指定时间间隔, 向前端发送ping消息心跳
	ticker := time.NewTicker(heartBeatPeriod)
	defer func() {
		ticker.Stop()
		c.close()
		if err := recover(); err != nil {
			global.Log.Error(c.ctx, "[消息中心][接收端][%s]连接可能已断开: %v, %s", c.Key, err, string(debug.Stack()))
		}
	}()
	for {
		select {
		// 到心跳检测时间
		case <-ticker.C:
			last := time.Now().Sub(c.LastActiveTime.Time)
			if c.RetryCount > HeartBeatMaxRetryCount {
				panic(fmt.Sprintf("尝试发送心跳多次(%d)无响应", c.RetryCount))
			}
			if last > heartBeatPeriod {
				// 发送心跳
				c.Send.SafeSend(response.MessageWsResp{
					Type:   MessageRespHeartBeat,
					Detail: resp.GetSuccessWithData(c.RetryCount),
				})
				c.RetryCount++
			}
		}
	}
}

// 用户上线
func (c *MessageClient) register() {
	hub.lock.Lock()
	defer hub.lock.Unlock()

	t := carbon.Now()
	active, ok := hub.UserLastActive[c.User.Id]
	last := carbon.CreateFromTimestamp(active)
	hub.Clients[c.Key] = c
	if !ok || last.AddDuration(lastActiveRegisterPeriod.String()).Lt(t) {
		if !utils.ContainsUint(hub.UserIds, c.User.Id) {
			hub.UserIds = append(hub.UserIds, c.User.Id)
		}
		global.Log.Debug(c.ctx, "[消息中心][用户上线][%s]%d-%s", c.Key, c.User.Id, c.Ip)
		go func() {
			hub.RefreshUserMessage.SafeSend([]uint{c.User.Id})
		}()

		// 广播当前用户上线
		msg := response.MessageWsResp{
			Type: MessageRespOnline,
			Detail: resp.GetSuccessWithData(map[string]interface{}{
				"user": c.User,
			}),
		}

		// 通知除自己之外的人
		go hub.Broadcast.SafeSend(MessageBroadcast{
			MessageWsResp: msg,
			UserIds:       utils.ContainsUintThenRemove(hub.UserIds, c.User.Id),
		})

		// 记录最后活跃时间戳
		hub.UserLastActive[c.User.Id] = t.Timestamp()
	} else {
		hub.UserLastActive[c.User.Id] = t.Timestamp()
	}
}

// 关闭连接
func (c *MessageClient) close() {
	hub.lock.Lock()
	defer hub.lock.Unlock()

	if _, ok := hub.Clients[c.Key]; ok {
		delete(hub.Clients, c.Key)
		// 关闭发送通道
		c.Send.SafeClose()
		global.Log.Debug(c.ctx, "[消息中心][用户下线][%s]%d-%s", c.Key, c.User.Id, c.Ip)
	}

	c.Conn.Close()
}
