package v1

import (
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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
	// 心跳消息
	MessageReqHeartBeat uint = 0
	// 连接断开
	MessageReqDisconnect uint = 1
	// 推送新消息
	MessageReqPush uint = 2
	// 批量已读
	MessageReqBatchRead uint = 3
	// 批量删除
	MessageReqBatchDeleted uint = 4
	// 全部已读
	MessageReqAllRead uint = 5
	// 全部删除
	MessageReqAllDeleted uint = 6

	// 消息响应类型
	// 心跳消息
	MessageRespHeartBeat uint = 0
	// 普通消息
	MessageRespNormal uint = 1
	// 连接断开
	MessageRespDisconnect uint = 2
	// 未读数
	MessageRespUnRead uint = 3
	// 用户上线
	MessageRespOnline uint = 4
	// 用户离线
	MessageRespOffline uint = 5
)

var (
	upgrade = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	hub MessageHub
)

// 消息仓库, 用于维护整个消息中心连接
type MessageHub struct {
	// mysql连接
	Mysql service.MysqlService
	// 客户端用户id集合
	UserIds []uint
	// 客户端集合(用户id为每个socket key)
	Clients map[string]*MessageClient
	// 客户端注册(用户上线通道)
	Register chan *MessageClient
	// 客户端取消注册(用户下线通道)
	UnRegister chan *MessageClient
	// 广播通道
	Broadcast chan MessageResp
	// 刷新用户消息通道
	RefreshUserMessage chan []uint
}

// 消息客户端
type MessageClient struct {
	// 当前socket key
	Key string
	// 当前socket连接实例
	Conn *websocket.Conn
	// 当前登录用户
	User models.SysUser
	// 发送消息通道
	Send chan MessageResp
	// 上次活跃时间
	LastActiveTime models.LocalTime
	// 重试次数
	RetryCount uint
}

// 消息请求
type MessageReq struct {
	// 消息类型, 见const
	Type uint `json:"type"`
	// 数据内容
	Data interface{} `json:"data"`
}

// 消息响应
type MessageResp struct {
	// 消息类型, 见const
	Type uint `json:"type"`
	// 消息详情
	Detail response.Resp `json:"detail"`
}

// 启动消息中心仓库
func StartMessageHub() {
	// 初始化
	hub.Mysql = service.New(nil)
	hub.Clients = make(map[string]*MessageClient)
	hub.Register = make(chan *MessageClient)
	hub.UnRegister = make(chan *MessageClient)
	hub.Broadcast = make(chan MessageResp)
	hub.RefreshUserMessage = make(chan []uint)
	go hub.run()
}

// 启动消息连接
func MessageWs(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Log.Error("创建消息连接失败", err)
		return
	}

	// 获取当前登录用户
	user := GetCurrentUser(c)

	// 注册到消息仓库
	client := &MessageClient{
		Key:  c.Request.Header.Get("Sec-WebSocket-Key"),
		Conn: conn,
		User: user,
		Send: make(chan MessageResp),
	}
	hub.Register <- client

	// 监听数据的接收/发送/心跳
	go client.receive()
	go client.send()
	go client.heartBeat()

	// 刷新用户消息
	hub.RefreshUserMessage <- []uint{user.Id}

	// 广播当前用户上线
	msg := MessageResp{
		Type: MessageRespOnline,
		Detail: response.GetSuccessWithData(map[string]interface{}{
			"user": user,
		}),
	}
	hub.Broadcast <- msg
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
			global.Log.Debug("用户上线: ", client.Key)
		// 用户下线	
		case client := <-h.UnRegister:
			if _, ok := h.Clients[client.Key]; ok {
				delete(h.Clients, client.Key)
				// 关闭发送通道
				close(client.Send)
				global.Log.Debug("用户下线: ", client.Key)
			}
		// 广播(全部用户均可接收)	
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				}
			}
		// 刷新客户端消息	
		case userIds := <-h.RefreshUserMessage:
			// 同步用户消息
			hub.Mysql.SyncMessageByUserIds(userIds)
			for _, client := range h.Clients {
				for _, id := range userIds {
					if client.User.Id == id {
						// 获取未读消息条数
						total, _ := hub.Mysql.GetUnReadMessageCount(id)
						// 将当前消息条数发送给用户
						msg := MessageResp{
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
	for {
		_, msg, err := c.Conn.ReadMessage()

		// 记录活跃时间
		c.LastActiveTime = models.LocalTime{
			Time: time.Now(),
		}

		if err != nil {
			global.Log.Error("接收数据失败, 连接可能已断开", err)
			hub.UnRegister <- c
			break
		}
		global.Log.Debug("接收数据成功", c.User.Id, string(msg))
		// 数据转为json
		var req MessageReq
		utils.Json2Struct(string(msg), &req)
		switch req.Type {
		case MessageReqHeartBeat:
			// 重置计数器
			c.LastActiveTime = models.LocalTime{
				Time: time.Now(),
			}
			c.RetryCount = 0
		case MessageReqPush:
			var data request.PushMessageRequestStruct
			utils.Struct2StructByJson(req.Data, &data)
			// 参数校验
			err = global.NewValidatorError(global.Validate.Struct(data), data.FieldTrans())
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			} else {
				data.FromUserId = c.User.Id
				err = hub.Mysql.CreateMessage(&data)
				if err != nil {
					detail = response.GetFailWithMsg(err.Error())
				}
			}
			if err == nil {
				// 刷新条数
				hub.RefreshUserMessage <- hub.UserIds
			}
			// 发送响应
			c.Send <- MessageResp{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqBatchRead:
			var data request.Req
			utils.Struct2StructByJson(req.Data, &data)
			err = hub.Mysql.BatchUpdateMessageRead(data.GetUintIds())
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- MessageResp{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqBatchDeleted:
			var data request.Req
			utils.Struct2StructByJson(req.Data, &data)
			err = hub.Mysql.BatchUpdateMessageDeleted(data.GetUintIds())
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- MessageResp{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqAllRead:
			err = hub.Mysql.UpdateAllMessageRead(c.User.Id)
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- MessageResp{
				Type:   MessageRespNormal,
				Detail: detail,
			}
		case MessageReqAllDeleted:
			err = hub.Mysql.UpdateAllMessageDeleted(c.User.Id)
			detail := response.GetSuccess()
			if err != nil {
				detail = response.GetFailWithMsg(err.Error())
			}
			// 刷新条数
			hub.RefreshUserMessage <- hub.UserIds
			// 发送响应
			c.Send <- MessageResp{
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
	}()
	for {
		select {
		// 发送通道
		case msg, ok := <-c.Send:
			// 设定回写超时时间
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// send通道已关闭
				c.Conn.WriteMessage(websocket.CloseMessage, []byte("closed"))
				// 强制下线
				hub.UnRegister <- c
				return
			}

			// 发送文本消息
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(utils.Struct2Json(msg))); err != nil {
				global.Log.Error("发送数据失败, 连接可能已断开", err)
				// 强制下线
				hub.UnRegister <- c
				return
			}
		// 长时间无新消息	
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 发送ping消息
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				global.Log.Error("发送数据失败, 连接可能已断开", err)
				// 强制下线
				hub.UnRegister <- c
				return
			}
		}
	}
}

// 心跳检测
func (c *MessageClient) heartBeat() {
	// 创建定时器, 超出指定时间间隔, 向前端发送ping消息心跳
	ticker := time.NewTicker(heartBeatPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		if err := recover(); err != nil {
			global.Log.Error("发送心跳失败, 连接可能已断开", err)
		}
	}()
loop:
	for {
		select {
		// 到心跳检测时间
		case <-ticker.C:
			global.Log.Debug("当前活跃连接", hub.Clients)
			last := time.Now().Sub(c.LastActiveTime.Time)
			if c.RetryCount > HeartBeatMaxRetryCount {
				global.Log.Error("尝试发送心跳多次无响应, 连接可能已断开")
				hub.UnRegister <- c
				break loop
			}
			if last > heartBeatPeriod {
				// 发送心跳
				c.Send <- MessageResp{
					Type:   MessageRespHeartBeat,
					Detail: response.GetSuccessWithData(c.RetryCount),
				}
				c.RetryCount++
			}
		}
	}
}
