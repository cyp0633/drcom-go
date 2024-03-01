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
			if server == "" {
				server = getServer()
			}
			if server != "" {
				doLogin(server)
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
