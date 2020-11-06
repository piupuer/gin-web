package utils

import (
	"fmt"
	"testing"
)

func TestExecRemoteShell(t *testing.T) {
	fmt.Println(ExecRemoteShell(SshConfig{
		LoginName: "gsgc",
		LoginPwd:  "123456",
		Host:      "192.168.5.108",
		Port:      22,
	}, []string{
		"ls -lh",
		"rm Pictures",
	}))
}
