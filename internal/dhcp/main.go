// 用于处理 DHCP 模式的登录与保活
package dhcp

import (
	"bufio"
	"net"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

var writer *bufio.Writer

// 使用 DHCP 模式上网
func Run() {
	util.Logger.Info("DHCP mode")
	c, err := net.Dial("udp", util.CLI.BindIP+":61440") // 绑定 udp/61440 端口与服务器通信
	if err != nil {
		util.Logger.Fatal("Open socket on 61440 failed", zap.Error(err))
	}
	// TODO: 超时动作，需要每次读写设置 deadline？
	writer = bufio.NewWriter(c)
	defer c.Close()
	for {
		tail, salt, err := login()
		if err != nil {
			switch err {
			case ErrorLogin:
				continue
			default:
				util.Logger.Fatal("Login failed", zap.Error(err))
			}
		}
		// 清除 socket buffer
		err = writer.Flush()
		if err != nil {
			util.Logger.Error("Flush socket failed", zap.Error(err))
		}
		// 保活
		keepAlive1(salt, tail)
		keepAlive2(salt, tail)
	}
}
