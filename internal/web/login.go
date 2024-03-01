// 用于模拟网页登录
package web

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

// 发送登录请求
func doLogin(server_base string) {
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
	util.Logger.Debug("Sending login request", zap.String("server", server_base), zap.String("form", form.Encode()))

	response, err := http.PostForm(server_base+"/0.htm", form)
	if err != nil {
		util.Logger.Error("Send login request failed", zap.Error(err))
	}
	util.Logger.Info("Login seems successful")
	util.Logger.Debug("Response", zap.Any("response", response.Body))
	defer response.Body.Close()
}
