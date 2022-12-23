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
	util.Logger.Debug("DHCP mode")
	laddr := net.UDPAddr{IP: net.ParseIP(util.CLI.BindIP), Port: 61440}
	raddr := net.UDPAddr{IP: net.ParseIP(util.Conf.Server), Port: 61440}
	c, err := net.DialUDP("udp", &laddr, &raddr)
	if err != nil {
		util.Logger.Fatal("Open socket on 61440 failed", zap.Error(err))
	}
	util.Logger.Debug("Opened socket on 61440", zap.String("local", c.LocalAddr().String()), zap.String("remote", c.RemoteAddr().String()))
	// TODO: 超时动作，需要每次读写设置 deadline？
	conn = bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	defer c.Close()
	for {
		tail, salt, err := login()
		if err != nil {
			time.Sleep(5 * time.Second)
			util.Logger.Info("Login failed, retrying", zap.Error(err))
			continue
		}
		// 清除 socket buffer
		err = conn.Flush()
		if err != nil {
			util.Logger.Error("Flush socket failed", zap.Error(err))
		}
		keepAlive2Counter = 0
		var first *int
		*first = 1
		// 保活
		for try := 0; try <= 5; {
			if err = keepAlive1(tail, salt); err == nil {
				time.Sleep(time.Microsecond * 200)
				err = keepAlive2(salt, tail)
				if err != nil {
					util.Logger.Info("Keepalive2 failed, retrying", zap.Error(err))
					time.Sleep(time.Second)
				} else {
					time.Sleep(time.Second * 20)
				}
			} else {
				try++
				util.Logger.Info("Keepalive1 failed, retrying", zap.Error(err))
				time.Sleep(time.Second)
			}
		}
	}
}
