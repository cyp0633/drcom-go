// 用于模拟网页登录
package web

import (
	"crypto/md5"
	"net/http"
	"net/url"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

// 发送登录请求
func doLogin(server_base string) {
	username := util.Conf.Username
	password := util.Conf.Password
	md5Ctx := md5.New()
	md5Ctx.Write([]byte("2" + password + "12345678"))
	cipherStr := md5Ctx.Sum(nil)

	form := url.Values{
		"DDDDD":  {username},
		"upass":  {string(cipherStr)},
		"0MKKey": {"123456"},
		"R1":     {"0"},
		"R2":     {"1"},
		"para":   {"00"},
		"v6ip":   {""},
	}
	util.Logger.Debug("Sending login request", zap.String("server", server_base),zap.String("form",form.Encode()))

	response, err := http.PostForm(server_base+"/0.htm", form)
	if err != nil {
		util.Logger.Error("Send login request failed", zap.Error(err))
	}
	util.Logger.Info("Login seems successful")
	util.Logger.Debug("Response", zap.Any("response", response.Body))
	defer response.Body.Close()
}
