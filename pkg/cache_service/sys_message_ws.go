package cache_service

import (
	"context"
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	"github.com/piupuer/go-helper/pkg/ch"
	"github.com/piupuer/go-helper/pkg/middleware"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Message websocket

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Send heartbeat to peer with this period
	heartBeatPeriod = 10 * time.Second

	// Online notifications will not be sent repeatedly during the cycle
	lastActiveRegisterPeriod = 10 * time.Minute

	// Maximum number of heartbeat retries
	heartBeatMaxRetryCount = 3

	// message type
	// first number: request/response
	// second number: message type
	// third number: same message type's sort
	// 
	// request message type(first number=1)
	MessageReqHeartBeat    string = "1-1-1"
	MessageReqPush         string = "1-2-1"
	MessageReqBatchRead    string = "1-2-2"
	MessageReqBatchDeleted string = "1-2-3"
	MessageReqAllRead      string = "1-2-4"
	MessageReqAllDeleted   string = "1-2-5"

	// response message type(first number=2)
	MessageRespHeartBeat string = "2-1-1"
	MessageRespNormal    string = "2-2-1"
	MessageRespUnRead    string = "2-3-1"
	MessageRespOnline    string = "2-4-1"
)

var hub MessageHub

// The message hub is used to maintain the entire message connection
type MessageHub struct {
	lock    sync.RWMutex
	Service RedisService
	// client user ids
	UserIds []uint
	// client user active timestamp
	UserLastActive map[uint]int64
	// MessageClients(key is websocket key)
	Clients map[string]*MessageClient
	// The channel sends messages to all users
	Broadcast *ch.Ch
	// The channel to refresh user message
	RefreshUserMessage *ch.Ch
	// Idempotence middleware options
	IdempotenceOps *middleware.IdempotenceOptions
}

// The message client is used to store connection information
type MessageClient struct {
	ctx context.Context
	// websocket key
	Key string
	// websocket connection
	Conn *websocket.Conn
	// current user
	User models.SysUser
	// request ip
	Ip string
	// The channel sends messages to user
	Send           *ch.Ch
	LastActiveTime carbon.Carbon
	RetryCount     uint
}

// The message broadcast is used to store users who need to broadcast
type MessageBroadcast struct {
	response.MessageWsResp
	UserIds []uint `json:"-"`
}

// start hub
func (rd RedisService) StartMessageHub(idempotenceOps *middleware.IdempotenceOptions) MessageHub {
	hub.Service = rd
	hub.Clients = make(map[string]*MessageClient)
	hub.UserLastActive = make(map[uint]int64)
	hub.Broadcast = ch.NewCh()
	hub.RefreshUserMessage = ch.NewCh()
	hub.IdempotenceOps = idempotenceOps
	go hub.run()
	go hub.count()
	return hub
}

// websocket handler
func (rd RedisService) MessageWs(ctx *gin.Context, conn *websocket.Conn, key string, user models.SysUser, ip string) {
	client := &MessageClient{
		ctx:  ctx,
		Key:  key,
		Conn: conn,
		User: user,
		Ip:   ip,
		Send: ch.NewCh(),
	}

	go client.register()
	go client.receive()
	go client.send()
	go client.heartBeat()
}

func (h MessageHub) run() {
	for {
		select {
		// Broadcast channel
		case data := <-h.Broadcast.C:
			broadcast := data.(MessageBroadcast)
			for _, client := range h.getClients() {
				// Notify the specified user
				if utils.ContainsUint(broadcast.UserIds, client.User.Id) {
					client.Send.SafeSend(broadcast)
				}
			}
			// RefreshUserMessage channel
		case data := <-h.RefreshUserMessage.C:
			userIds := data.([]uint)
			// sync users message
			hub.Service.mysql.SyncMessageByUserIds(userIds)
			for _, client := range h.getClients() {
				for _, id := range userIds {
					if client.User.Id == id {
						total, _ := hub.Service.mysql.GetUnReadMessageCount(id)
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

// active connection count
func (h MessageHub) count() {
	ticker := time.NewTicker(heartBeatPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			infos := make([]string, 0)
			for _, client := range h.getClients() {
				infos = append(infos, fmt.Sprintf("%d-%s", client.User.Id, client.Ip))
			}
			global.Log.Debug(h.Service.Q.Ctx, "[Message]active connection: %v", strings.Join(infos, ","))
		}
	}
}

func (h MessageHub) getClients() map[string]*MessageClient {
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	return hub.Clients
}

// receive message handler
func (c *MessageClient) receive() {
	defer func() {
		c.close()
		if err := recover(); err != nil {
			global.Log.Error(c.ctx, "[Message][receiver][%s]connection may have been lost: %v, %s", c.Key, err, string(debug.Stack()))
		}
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()

		// save active time
		c.LastActiveTime = carbon.Now()
		c.RetryCount = 0

		if err != nil {
			panic(err)
		}
		// decompress data
		// data := utils.DeCompressStrByZlib(string(msg))
		data := string(msg)
		global.Log.Debug(c.ctx, "[Message][receiver][%s]receive data success: %d, %s", c.Key, c.User.Id, data)
		var r request.MessageWsReq
		utils.Json2Struct(data, &r)
		switch r.Type {
		case MessageReqHeartBeat:
			if _, ok := r.Data.(float64); ok {
				c.Send.SafeSend(response.MessageWsResp{
					Type:   MessageRespHeartBeat,
					Detail: resp.GetSuccess(),
				})
			}
		case MessageReqPush:
			var data request.PushMessageReq
			utils.Struct2StructByJson(r.Data, &data)
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
				hub.RefreshUserMessage.SafeSend(hub.UserIds)
			}
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
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
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
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
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
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
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
			hub.RefreshUserMessage.SafeSend(hub.UserIds)
			c.Send.SafeSend(response.MessageWsResp{
				Type:   MessageRespNormal,
				Detail: detail,
			})
		}
	}
}

// send message handler
func (c *MessageClient) send() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
		if err := recover(); err != nil {
			global.Log.Error(c.ctx, "[Message][sender][%s]connection may have been lost: %v, %s", c.Key, err, string(debug.Stack()))
		}
	}()
	for {
		select {
		case msg, ok := <-c.Send.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// send failed
				c.writeMessage(websocket.CloseMessage, "closed")
				panic("connection closed")
			}

			if err := c.writeMessage(websocket.TextMessage, utils.Struct2Json(msg)); err != nil {
				panic(err)
			}
		// timeout
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.writeMessage(websocket.PingMessage, "ping"); err != nil {
				panic(err)
			}
		}
	}
}

func (c MessageClient) writeMessage(messageType int, data string) error {
	// compress
	// s, _ := utils.CompressStrByZlib(data)
	s := &data
	global.Log.Debug(c.ctx, "[Message][sender][%s] %v", c.Key, *s)
	return c.Conn.WriteMessage(messageType, []byte(*s))
}

// heartbeat handler
func (c *MessageClient) heartBeat() {
	ticker := time.NewTicker(heartBeatPeriod)
	defer func() {
		ticker.Stop()
		c.close()
		if err := recover(); err != nil {
			global.Log.Error(c.ctx, "[Message][heartbeat][%s]connection may have been lost: %v, %s", c.Key, err, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-ticker.C:
			last := time.Now().Sub(c.LastActiveTime.Time)
			if c.RetryCount > heartBeatMaxRetryCount {
				panic(fmt.Sprintf("[Message][heartbeat]retry sending heartbeat for %d times without response", c.RetryCount))
			}
			if last > heartBeatPeriod {
				c.Send.SafeSend(response.MessageWsResp{
					Type:   MessageRespHeartBeat,
					Detail: resp.GetSuccessWithData(c.RetryCount),
				})
				c.RetryCount++
			} else {
				c.RetryCount = 0
			}
		}
	}
}

// user online handler
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
		global.Log.Debug(c.ctx, "[Message][online][%s]%d-%s", c.Key, c.User.Id, c.Ip)
		go func() {
			hub.RefreshUserMessage.SafeSend([]uint{c.User.Id})
		}()

		msg := response.MessageWsResp{
			Type: MessageRespOnline,
			Detail: resp.GetSuccessWithData(map[string]interface{}{
				"user": c.User,
			}),
		}
		// Inform everyone except yourself
		go hub.Broadcast.SafeSend(MessageBroadcast{
			MessageWsResp: msg,
			UserIds:       utils.ContainsUintThenRemove(hub.UserIds, c.User.Id),
		})

		hub.UserLastActive[c.User.Id] = t.Timestamp()
	} else {
		hub.UserLastActive[c.User.Id] = t.Timestamp()
	}
}

// user offline handler
func (c *MessageClient) close() {
	hub.lock.Lock()
	defer hub.lock.Unlock()

	if _, ok := hub.Clients[c.Key]; ok {
		delete(hub.Clients, c.Key)
		c.Send.SafeClose()
		global.Log.Debug(c.ctx, "[Message][offline][%s]%d-%s", c.Key, c.User.Id, c.Ip)
	}

	c.Conn.Close()
}
