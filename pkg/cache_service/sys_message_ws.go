package cache_service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gorilla/websocket"
	"strings"
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
	// redis连接
	Service RedisService
	// 客户端用户id集合
	UserIds []uint
	// 客户端集合(用户id为每个socket key)
	Clients map[string]*MessageClient
	// 客户端注册(用户上线通道)
	Register chan *MessageClient
	// 客户端取消注册(用户下线通道)
	UnRegister chan *MessageClient
	// 广播通道
	Broadcast chan MessageBroadcast
	// 刷新用户消息通道
	RefreshUserMessage chan []uint
	// 幂等性token校验方法
	CheckIdempotenceTokenFunc func(token string) bool
}

// 消息客户端
type MessageClient struct {
	// 当前socket key
	Key string
	// 当前socket连接实例
	Conn *websocket.Conn
	// 当前登录用户
	User models.SysUser
	// 当前登录用户ip地址
	Ip string
	// 发送消息通道
	Send chan response.MessageWsResponseStruct
	// 上次活跃时间
	LastActiveTime models.LocalTime
	// 重试次数
	RetryCount uint
}

// 消息广播
type MessageBroadcast struct {
	response.MessageWsResponseStruct
	UserIds []uint `json:"-"`
}

// 启动消息中心仓库
func (s *RedisService) StartMessageHub(checkIdempotenceTokenFunc func(token string) bool) MessageHub {
	// 初始化参数
	hub.Service = *s
	hub.Clients = make(map[string]*MessageClient)
	hub.Register = make(chan *MessageClient)
	hub.UnRegister = make(chan *MessageClient)
	hub.Broadcast = make(chan MessageBroadcast)
	hub.RefreshUserMessage = make(chan []uint)
	hub.CheckIdempotenceTokenFunc = checkIdempotenceTokenFunc
	go hub.run()
	return hub
}

// 启动消息连接
func (s *RedisService) MessageWs(conn *websocket.Conn, key string, user models.SysUser, ip string) {
	// 注册到消息仓库
	client := &MessageClient{
		Key:  key,
		Conn: conn,
		User: user,
		Ip:   ip,
		Send: make(chan response.MessageWsResponseStruct),
	}
	hub.Register <- client

	// 监听数据的接收/发送/心跳
	go client.receive()
	go client.send()
	// go client.heartBeat()

	// 刷新用户消息
	hub.RefreshUserMessage <- []uint{user.Id}

	// 广播当前用户上线
	msg := response.MessageWsResponseStruct{
		Type: MessageRespOnline,
		Detail: response.GetSuccessWithData(map[string]interface{}{
			"user": user,
		}),
	}

	// 通知除自己之外的人
	hub.Broadcast <- MessageBroadcast{
		MessageWsResponseStruct: msg,
		UserIds:                 utils.ContainsUintThenRemove(hub.UserIds, user.Id),
	}
}

// 运行仓库
func (h *MessageHub) run() {
	for {
		select {
		// 新用户上线
		case client := <-h.Register:
			if !utils.ContainsUint(h.UserIds, client.User.Id) {
				h.UserIds = append(h.UserIds, client.User.Id)
			}
			h.Clients[client.Key] = client
			global.Log.Debug("[消息中心][广播]用户上线: ", fmt.Sprintf("%d-%s", client.User.Id, client.Ip))
		// 用户下线
		case client := <-h.UnRegister:
			if _, ok := h.Clients[client.Key]; ok {
				delete(h.Clients, client.Key)
				// 关闭发送通道
				close(client.Send)
				global.Log.Debug("[消息中心][广播]用户下线: ", fmt.Sprintf("%d-%s", client.User.Id, client.Ip))
			}
		// 广播(全部用户均可接收)
		case broadcast := <-h.Broadcast:
			for _, client := range h.Clients {
				// 通知指定用户
				if utils.ContainsUint(broadcast.UserIds, client.User.Id) {
					select {
					case client.Send <- broadcast.MessageWsResponseStruct:
					}
				}
			}
		// 刷新客户端消息
		case userIds := <-h.RefreshUserMessage:
			// 同步用户消息
			hub.Service.mysql.SyncMessageByUserIds(userIds)
			for _, client := range h.Clients {
				for _, id := range userIds {
					if client.User.Id == id {
						// 获取未读消息条数
						total, _ := hub.Service.mysql.GetUnReadMessageCount(id)
						// 将当前消息条数发送给用户
						msg := response.MessageWsResponseStruct{
							Type: MessageRespUnRead,
							Detail: response.GetSuccessWithData(map[string]int64{
								"unReadCount": total,
							}),
						}
						client.Send <- msg
					}
				}
			}
		}
	}
}

// 接收数据
func (c *MessageClient) receive() {
	defer func() {
		c.Conn.Close()
		if err := recover(); err != nil {
			global.Log.Error("[消息中心][接收端]连接可能已断开: ", err)
		}
	}()
loop:
	for {
		_, msg, err := c.Conn.ReadMessage()

		// 记录活跃时间
		c.LastActiveTime = models.LocalTime{
			Time: time.Now(),
		}
		c.RetryCount = 0

		if err != nil {
			global.Log.Error("[消息中心][接收端]接收数据失败: ", err)
			hub.UnRegister <- c
			break loop
		}
		// 解压数据
		// data := utils.DeCompressStrByZlib(string(msg))
		data := string(msg)
		global.Log.Debug("[消息中心][接收端]接收数据成功: ", c.User.Id, data)
		// 数据转为json
		var req request.MessageWsRequestStruct
		utils.Json2Struct(data, &req)
		switch req.Type {
		case MessageReqHeartBeat:
			if _, ok := req.Data.(float64); ok {
				// 发送心跳
				c.Send <- response.MessageWsResponseStruct{
					Type:   MessageRespHeartBeat,
					Detail: response.GetSuccess(),
				}
			}
		case MessageReqPush:
			var data request.PushMessageRequestStruct
			utils.Struct2StructByJson(req.Data, &data)
			// 参数校验
			err = global.NewValidatorError(global.Validate.Struct(data), data.FieldTrans())
			detail := response.GetSuccess()
			if err == nil {
				if !hub.CheckIdempotenceTokenFunc(data.IdempotenceToken) {
					err = errors.New(response.IdempotenceTokenInvalidMsg)
				} else {
					data.FromUserId = c.User.Id
					err = hub.Service.mysql.CreateMessage(&data)
				}
			}
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			} else {
				// 刷新条数
				hub.RefreshUserMessage <- hub.UserIds
			}
			// 发送响应
			c.Send <- response.MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqBatchRead:
			var data request.Req
			utils.Struct2StructByJson(req.Data, &data)
			err = hub.Service.mysql.BatchUpdateMessageRead(data.GetUintIds())
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- response.MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqBatchDeleted:
			var data request.Req
			utils.Struct2StructByJson(req.Data, &data)
			err = hub.Service.mysql.BatchUpdateMessageDeleted(data.GetUintIds())
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- response.MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqAllRead:
			err = hub.Service.mysql.UpdateAllMessageRead(c.User.Id)
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- response.MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqAllDeleted:
			err = hub.Service.mysql.UpdateAllMessageDeleted(c.User.Id)
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- response.MessageWsResponseStruct{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		}
	}
}

// 发送数据
func (c *MessageClient) send() {
	// 创建定时器, 超出指定时间间隔, 向前端发送ping消息心跳
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		if err := recover(); err != nil {
			global.Log.Error("[消息中心][发送端]连接可能已断开: ", err)
		}
	}()
	for {
		select {
		// 发送通道
		case msg, ok := <-c.Send:
			// 设定回写超时时间
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// send通道已关闭
				c.writeMessage(websocket.CloseMessage, "closed")
				// 强制下线
				hub.UnRegister <- c
				return
			}

			// 发送文本消息
			if err := c.writeMessage(websocket.TextMessage, utils.Struct2Json(msg)); err != nil {
				global.Log.Error("[消息中心][发送端]发送数据失败: ", err)
				// 强制下线
				hub.UnRegister <- c
				return
			}
		// 长时间无新消息
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 发送ping消息
			if err := c.writeMessage(websocket.PingMessage, "ping"); err != nil {
				global.Log.Error("[消息中心][发送端]发送数据失败: ", err)
				// 强制下线
				hub.UnRegister <- c
				return
			}
		}
	}
}

// 回写消息
func (c *MessageClient) writeMessage(messageType int, data string) error {
	// 字符串压缩
	// s, _ := utils.CompressStrByZlib(data)
	s := &data
	global.Log.Debug("[消息中心][writeMessage]", *s)
	return c.Conn.WriteMessage(messageType, []byte(*s))
}

// 心跳检测
func (c *MessageClient) heartBeat() {
	// 创建定时器, 超出指定时间间隔, 向前端发送ping消息心跳
	ticker := time.NewTicker(heartBeatPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		if err := recover(); err != nil {
			global.Log.Error("[消息中心][心跳]连接可能已断开: ", err)
		}
	}()
loop:
	for {
		select {
		// 到心跳检测时间
		case <-ticker.C:
			infos := make([]string, 0)
			for _, client := range hub.Clients {
				infos = append(infos, fmt.Sprintf("%d-%s", client.User.Id, client.Ip))
			}
			global.Log.Debug("[消息中心][心跳]当前活跃连接: ", strings.Join(infos, ","))
			last := time.Now().Sub(c.LastActiveTime.Time)
			if c.RetryCount > HeartBeatMaxRetryCount {
				global.Log.Error(fmt.Sprintf("[消息中心][心跳]尝试发送心跳多次(%d)无响应", c.RetryCount))
				hub.UnRegister <- c
				break loop
			}
			if last > heartBeatPeriod {
				// 发送心跳
				c.Send <- response.MessageWsResponseStruct{
					Type:   MessageRespHeartBeat,
					Detail: response.GetSuccessWithData(c.RetryCount),
				}
				c.RetryCount++
			}
		}
	}
}
