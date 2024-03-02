package web

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyp0633/drcom-go/internal/util"
)

// 认证服务器地址
var serverBase string

// 连接检查与保活
func Run() {
	// 捕捉退出，自动注销
	if util.CLI.AutoLogout {
		exitChan := make(chan os.Signal, 1)
		signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)
		go func(ch chan os.Signal) {
			<-ch
			doLogout()
		}(exitChan)
	}
	// 检查连接，设置通信
	ch := make(chan bool, 1)
	go util.CheckConnection(ch, 1*time.Minute) // 似乎比较不容易断，可以设长一点
	for {
		select {
		case <-ch:
			util.Logger.Info("Network disconnected, trying to login")
			retry := 0
			for serverBase == "" { // 时不时获取不到解析结果的 workaround；有网则不会进入此循环
				if retry != 0 {
					time.Sleep(5 * time.Second)
				}
				if retry >= 5 {
					util.Logger.Error("Get captive server failed too many times")
					break
				}
				serverBase = getServer()
				retry++
			}
			if serverBase != "" {
				doLogin()
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
