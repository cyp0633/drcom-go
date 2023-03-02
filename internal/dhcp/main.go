// 用于处理 DHCP 模式的登录与保活
package dhcp

import (
	"bufio"
	"math/rand"
	"net"
	"time"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

var conn *bufio.ReadWriter

var udpConn *net.UDPConn

// 使用 DHCP 模式上网
func Run() {
	util.Logger.Debug("DHCP mode")
	laddr := net.UDPAddr{IP: net.ParseIP(util.CLI.BindIP), Port: 61440}
	raddr := net.UDPAddr{IP: net.ParseIP(util.Conf.Server), Port: 61440}
	c, err := net.DialUDP("udp", &laddr, &raddr)
	udpConn = c
	if err != nil {
		util.Logger.Fatal("Open socket on 61440 failed", zap.Error(err))
	}
	util.Logger.Debug("Opened socket on 61440", zap.String("local", c.LocalAddr().String()), zap.String("remote", c.RemoteAddr().String()))
	// TODO: 超时动作，需要每次读写设置 deadline？
	conn = bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	defer c.Close()
login:
	for fail := 0; ; {
		// 不开启无限重试，且重试次数超过 5 次，则退出
		if !util.CLI.Eternal && fail > 5 {
			break
		}

		// 登录
		tail, salt, err := login()
		if err != nil {
			// 在随机的 [1,2^fail]*5 秒后重试（搁着指数退避是吧）
			var sleepTime = time.Second * time.Duration((rand.Intn((1<<fail)-1)+1)*5)
			util.Logger.Info("Login failed, retrying", zap.Error(err), zap.Duration("sleep", sleepTime), zap.Int("fail", fail))
			time.Sleep(sleepTime)
			fail++
			continue login
		} else {
			// 清除连续登录失败次数
			fail = 0
		}

		// 启动连接测试
		var ch = make(chan bool, 1)
		go util.CheckConnection(ch)

		time.Sleep(3 * time.Second)
		// 清除 socket buffer
		err = conn.Flush()
		if err != nil {
			util.Logger.Error("Flush socket failed", zap.Error(err))
		}
		var first = new(bool)
		*first = true

		// 保活
		for try := 0; try <= 5; {
			if err = keepAlive1(salt, tail); err == nil {
				time.Sleep(time.Microsecond * 200)
				err = keepAlive2(first, 0)
				if err != nil {
					util.Logger.Info("Keepalive2 failed, retrying", zap.Error(err))
					time.Sleep(time.Second)
				} else {
					// 如果 20 秒内没有检测到网络断开，则继续保活；否则重新登录
					select {
					case <-ch:
						util.Logger.Info("Recovering connection")
						continue login
					case <-time.After(time.Second * 20):
					}
				}
			} else {
				try++
				util.Logger.Info("Keepalive1 failed, retrying", zap.Error(err))
				time.Sleep(time.Second)
			}
		}
	}
	util.Logger.Error("Login failed 5 times, exiting")
}
