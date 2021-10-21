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

// start shell websocket
func MachineShellWs(c *gin.Context) {
	var req request.MachineShellWsReq
	err := c.ShouldBind(&req)

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Log.Error(c, "upgrade websocket failed: %v", err)
		return
	}
	defer conn.Close()

	active := time.Now()

	// get ssh client
	client, err := utils.GetSshClient(utils.SshConfig{
		Host:      req.Host,
		Port:      int(req.SshPort),
		LoginName: req.LoginName,
		LoginPwd:  req.LoginPwd,
	})
	if err != nil {
		global.Log.Error(c, "connect ssh failed: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	// open ssh channel
	channel, incomingRequests, err := client.Conn.OpenChannel("session", nil)
	if err != nil {
		global.Log.Error(c, "connect ssh failed: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}
	defer channel.Close()
	defer client.Close()

	go func() {
		for req := range incomingRequests {
			if req.WantReply {
				// reply 
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

	rows := uint32(req.Rows)
	cols := uint32(req.Cols)
	ptyReq := ptyRequestMsg{
		Term:     "xterm",
		Columns:  rows,
		Rows:     cols,
		Width:    rows,
		Height:   cols,
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&ptyReq))
	if !ok || err != nil {
		global.Log.Error(c, "send pseudo terminal request failed: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		global.Log.Error(c, "send shell failed: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("\n"+err.Error()))
		return
	}

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
					global.Log.Warn(c, "read shell failed: %v", err)
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
						global.Log.Error(c, "write msg to %s failed: %v", conn.RemoteAddr(), err)
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

	// timeout handler
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			cost := time.Since(active)
			if cost.Minutes() >= 30 {
				// if it is inactive for more than 30 minutes, the connection will be closed automatically
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\nif it is inactive for more than [ %s ] minutes, the connection will be closed automatically", cost.String())))
				conn.Close()
				break
			}
		}
	}()
	conn.WriteMessage(websocket.TextMessage, []byte("\r\n["+req.Host+"]connect success"))
	conn.WriteMessage(websocket.TextMessage, []byte("\r\nif it is inactive for more than 30 minutes, the connection will be closed automatically\r\n\r\n"))

	if req.InitCmd != "" {
		go func() {
			time.Sleep(time.Second * 1)
			channel.Write([]byte(req.InitCmd + "\r\n"))
		}()
	}

	// read user input
	for {
		m, p, err := conn.ReadMessage()
		active = time.Now()
		if err != nil {
			global.Log.Warn(c, "connect %s lost", conn.RemoteAddr())
			break
		}

		if m == websocket.TextMessage {
			cmd := string(p)
			if err := utils.IsSafetyCmd(cmd); err != nil {
				err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\n\r\n%s\r\n\r\n", err.Error())))
				if err != nil {
					global.Log.Warn(c, "write msg failed: %v", err)
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
