package web

import (
	"time"

	"github.com/cyp0633/drcom-go/internal/util"
)

// 连接检查与保活
// TODO: 退出时自动注销
func Run() {
	var server string
	ch := make(chan bool, 1)
	go util.CheckConnection(ch, 1*time.Minute) // 似乎比较不容易断，可以设长一点
	for {
		select {
		case <-ch:
			util.Logger.Info("Network disconnected, trying to login")
			retry := 0
			for server == "" { // 时不时获取不到解析结果的 workaround；有网则不会进入此循环
				if retry != 0 {
					time.Sleep(5 * time.Second)
				}
				if retry >= 5 {
					util.Logger.Error("Get captive server failed too many times")
					break
				}
				server = getServer()
				retry++
			}
			if server != "" {
				doLogin(server)
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
