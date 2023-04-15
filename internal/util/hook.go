package util

import (
	"os/exec"
	"strconv"
	"strings"
)

// 连接成功后执行的外部动作
func HookConnectSuccess() {
	exec.Command(ExtConf.ExecOnConnect).Start()
}

// 连接失败后执行的外部动作
//
// failureCount: 失败次数
func HookDisconnect(failureCount int) {
	cmd := ExtConf.ExecOnDisconnect
	// 将 {{FailureCount}} 替换为失败次数
	if cmd != "" {
		cmd = strings.Replace(cmd, "{{FailureCount}}", strconv.Itoa(failureCount), -1)
	}
	exec.Command(cmd).Start()
}
