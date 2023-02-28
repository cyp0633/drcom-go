package util

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// CheckConnection 检查网络连接，如发现不通，向 ch 中发送信号
func CheckConnection(ch chan bool) {
	for {
		resp, err := http.Get(ExtConf.KeepaliveServer)
		if err != nil || resp.StatusCode != http.StatusNoContent {
			ch <- true
			Logger.Warn("Network connection lost", zap.Error(err))
			return
		} else {
			resp.Body.Close()
			Logger.Debug("Network connection is OK")
			time.Sleep(time.Second * 5)
		}
	}
}
