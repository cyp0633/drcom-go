// 用于处理 DHCP 模式的登录与保活
package dhcp

import (
	"bufio"
	"net"
	"time"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

var conn *bufio.ReadWriter

// 使用 DHCP 模式上网
func Run() {
	util.Logger.Info("DHCP mode")
	laddr := net.UDPAddr{IP: net.ParseIP(util.CLI.BindIP), Port: 61440}
	raddr := net.UDPAddr{IP: net.ParseIP(util.Conf.Server), Port: 61440}
	c, err := net.DialUDP("udp", &laddr, &raddr)
	if err != nil {
		util.Logger.Fatal("Open socket on 61440 failed", zap.Error(err))
	}
	util.Logger.Info("Opened socket on 61440", zap.String("local", c.LocalAddr().String()), zap.String("remote", c.RemoteAddr().String()))
	// TODO: 超时动作，需要每次读写设置 deadline？
	conn = bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	defer c.Close()
	for {
		tail, salt, err := login()
		if err != nil {
			switch err {
			case ErrorLogin:
				time.Sleep(5 * time.Second)
				continue
			default:
				util.Logger.Fatal("Login failed", zap.Error(err))
			}
		}
		// 清除 socket buffer
		err = conn.Flush()
		if err != nil {
			util.Logger.Error("Flush socket failed", zap.Error(err))
		}
		// 保活
		keepAlive1(salt, tail)
		keepAlive2(salt, tail)
	}
}
