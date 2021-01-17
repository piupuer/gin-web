package utils

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"path"
	"strings"
	"time"
)

// ssh连接超时默认值
const DefaultSshTimeout = 5

type SshConfig struct {
	LoginName string
	LoginPwd  string
	Host      string
	Port      int
	Timeout   int
}

type SshResult struct {
	Connect bool   `json:"connect"`
	Result  string `json:"result"`
	Err     error  `json:"err"`
}

// 判断命令是否是允许的安全命令
func IsSafetyCmd(cmd string) error {
	// 避免rm * 或 rm /*等命令直接出现, 删除命令指定全路径
	c := path.Clean(strings.ToLower(cmd))
	if strings.Contains(c, "rm") {
		if len(strings.Split(c, "/")) <= 1 {
			return fmt.Errorf("rm命令%s不能删除小于2级目录的文件", cmd)
		}
	}
	return nil
}

// 执行远程命令
func ExecRemoteShell(config SshConfig, cmds []string) SshResult {
	return ExecRemoteShellWithTimeout(config, cmds, 0)
}

// 执行远程命令，超时自动关闭session，sessionCloseSecond=0 代表不设置超时强制关闭session
func ExecRemoteShellWithTimeout(config SshConfig, cmds []string, timeout int64) SshResult {
	var session *ssh.Session
	client, err := GetSshClient(config)
	if err != nil {
		return SshResult{
			Connect: false,
			Err:     err,
		}
	}
	// create session
	if session, err = client.NewSession(); err != nil {
		return SshResult{
			Connect: false,
			Err:     fmt.Errorf("无法创建ssh会话, 地址%s, 错误信息%v", config.Host, err),
		}
	}
	defer closeClient(session, client)

	go func() {
		if timeout > 0 {
			sleep, err := time.ParseDuration(fmt.Sprintf("%ds", timeout))
			if err != nil {
				fmt.Println(fmt.Sprintf("无法关闭ssh会话: %v", err))
				return
			}
			time.Sleep(sleep)
			closeClient(session, client)
		}
	}()

	command := ""

	for i, cmd := range cmds {
		if err := IsSafetyCmd(cmd); err != nil {
			return SshResult{
				Connect: true,
				Err:     err,
			}
		}
		if i == 0 {
			command = cmd
		} else {
			command = command + " && " + cmd
		}
	}

	var e bytes.Buffer
	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &e
	if command != "" {
		if err := session.Run(command); err != nil {
			return SshResult{
				Connect: true,
				Err:     fmt.Errorf("命令%s执行失败: %v", command, err),
				Result:  e.String(),
			}
		}
		fmt.Println(fmt.Sprintf("执行命令%s", command))
	}
	return SshResult{
		Result:  b.String(),
		Err:     nil,
		Connect: true,
	}
}

// 关闭连接
func closeClient(session *ssh.Session, client *ssh.Client) {
	err := client.Close()
	if err != nil {
		fmt.Println(fmt.Sprintf("关闭连接失败: %v", err))
	}
	// 必须关闭Client，才能释放该ssh连接句柄
	session.Close()
}

// 获取ssh连接
func GetSshClient(config SshConfig) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(config.LoginPwd))

	// 设定连接默认超时时间
	if config.Timeout == 0 {
		config.Timeout = DefaultSshTimeout
	}
	clientConfig = &ssh.ClientConfig{
		User:    config.LoginName,
		Auth:    auth,
		Timeout: time.Second * time.Duration(config.Timeout),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, fmt.Errorf("无法连接ssh, 地址%s, 错误信息%v", addr, err)
	}
	return client, nil
}
