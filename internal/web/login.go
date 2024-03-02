// 用于模拟网页登录
package web

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"os"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

// 发送登录请求
func doLogin() {
	username := util.Conf.Username
	password := util.Conf.Password
	hash := md5.Sum([]byte("2" + password + "12345678"))
	cipherStr := hex.EncodeToString(hash[:]) + "123456782"
	form := url.Values{
		"DDDDD":  {username},
		"upass":  {cipherStr},
		"0MKKey": {"123456"},
		"R1":     {"0"},
		"R2":     {"1"},
		"para":   {"00"},
		"v6ip":   {""},
	}
	util.Logger.Debug("Sending login request", zap.String("server", serverBase), zap.String("form", form.Encode()))

	response, err := http.PostForm(serverBase+"/0.htm", form)
	if err != nil {
		util.Logger.Error("Send login request failed", zap.Error(err))
	}
	util.Logger.Info("Login seems successful")
	util.Logger.Debug("Response", zap.Any("response", response.Body))
	defer response.Body.Close()
}

func doLogout() {
	if serverBase == "" { // 不是本程序登录的，不注销
		util.Logger.Warn("Captive server not set, not logging out...")
		return
	}
	// 退出：GET $serverBase/F.htm
	resp, err := http.Get(serverBase + "/F.htm")
	if err != nil {
		util.Logger.Error("Logout failed", zap.Error(err))
	} else {
		util.Logger.Info("Logout successful, exiting", zap.String("status", resp.Status))
	}
	os.Exit(0)
}
