package v1

import (
	"bufio"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"time"
	"unicode/utf8"
)

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

// 启动shell连接
func MachineShellWs(c *gin.Context) {
	var req request.MachineShellWsRequestStruct
	err := c.ShouldBind(&req)

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Log.Error("升级websocket连接失败", err)
		return
	}
	// 结束后自动关闭
	defer conn.Close()

	active := time.Now()

	// 建立连接
	client, err := utils.GetSshClient(utils.SshConfig{
		Host:      req.Host,
		Port:      req.SshPort,
		LoginName: req.LoginName,
		LoginPwd:  req.LoginPwd,
	})
	if err != nil {
		global.Log.Error(fmt.Sprintf("建立ssh连接失败%v", err))
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// 开启ssh通道
	channel, incomingRequests, err := client.Conn.OpenChannel("session", nil)
	if err != nil {
		global.Log.Error(fmt.Sprintf("建立ssh通道失败%v", err))
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}
	defer channel.Close()
	defer client.Close()

	// 处理需要回复的请求
	go func() {
		for req := range incomingRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}

	modeList = append(modeList, 0)

	// 发送pty
	ptyReq := ptyRequestMsg{
		Term:     "xterm",
		Columns:  req.Cols,
		Rows:     req.Rows,
		Width:    req.Cols * 8,
		Height:   req.Rows * 8,
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&ptyReq))
	if !ok || err != nil {
		global.Log.Error(fmt.Sprintf("发送pty失败%v", err))
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// 发送shell
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		global.Log.Error(fmt.Sprintf("发送shell失败%v", err))
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// 处理数据读写
	go func() {
		br := bufio.NewReader(channel)
		var buf []byte

		t := time.NewTimer(time.Millisecond * 100)
		defer t.Stop()
		r := make(chan rune)

		go func() {
			for {
				x, size, err := br.ReadRune()
				if err != nil {
					global.Log.Warn(fmt.Sprintf("读取shell警告%v", err))
					break
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		for {
			select {
			case <-t.C:
				if len(buf) != 0 {
					err = conn.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						global.Log.Error(fmt.Sprintf("数据写出到%s失败%v", conn.RemoteAddr(), err))
						return
					}
				}
				t.Reset(time.Millisecond * 100)
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}

	}()

	// 超时处理
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			cost := time.Since(active)
			if cost.Minutes() >= 30 {
				// 超过30分钟未活动，自动关闭连接
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\n已超过【%s】未活动，自动断开连接", cost.String())))
				conn.Close()
				break
			}
		}
	}()
	conn.WriteMessage(websocket.TextMessage, []byte("\r\n["+req.Host+"]终端服务连接成功"))
	conn.WriteMessage(websocket.TextMessage, []byte("\r\n提示：超过30分钟未活动，将自动断开连接\r\n\r\n"))

	if req.InitCmd != "" {
		go func() {
			time.Sleep(time.Second * 1)
			channel.Write([]byte(req.InitCmd + "\r\n"))
		}()
	}

	// 持续读取用户输入的命令
	for {
		m, p, err := conn.ReadMessage()
		active = time.Now()
		if err != nil {
			global.Log.Warn(fmt.Sprintf("连接%s已断开", conn.RemoteAddr()))
			break
		}

		if m == websocket.TextMessage {
			cmd := string(p)
			if err := utils.IsSafetyCmd(cmd); err != nil {
				err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\n\r\n%s\r\n\r\n", err.Error())))
				if err != nil {
					global.Log.Warn(fmt.Sprintf("回写数据失败%v", err))
					break
				}
				continue
			}
			if _, err := channel.Write(p); nil != err {
				break
			}
		}
	}
}
