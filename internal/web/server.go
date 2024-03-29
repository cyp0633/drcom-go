package web

import (
	"io"
	"net/http"
	"strings"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

// 查找登录服务器地址
func getServer() (server string) {
	generate_204 := util.ExtConf.ConnectionTestServer
	// 发 GET，检查是否有登录信息
	response, err := http.Get(generate_204)
	if err != nil {
		util.Logger.Error("Get server failed", zap.Error(err))
		return
	}
	defer response.Body.Close()
	util.Logger.Debug("Find auth server", zap.Any("response_header", response.Header))
	switch response.StatusCode {
	case http.StatusNoContent: // 正常访问，会返回 204
		util.Logger.Warn("204'ed; may have already logged in")
		return
	case http.StatusFound: // 临时重定向，Location header 即为认证服务器
		server = response.Header.Get("Location")
		if server == "" {
			util.Logger.Error("Server not found", zap.Any("headers", response.Header))
		} else {
			util.Logger.Info("Auth server found", zap.String("server", server))
		}
	case http.StatusOK: // 劫持了内容，但解析似乎没问题
		// 截取 location.href="<server>"
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			util.Logger.Error("Read body failed", zap.Error(err))
			return
		}
		bodyString := string(bodyBytes)
		startIndex := strings.Index(bodyString, "location.href=\"")
		if startIndex != -1 {
			startIndex += len("location.href=\"")
			endIndex := strings.Index(bodyString[startIndex:], "\"")
			if endIndex != -1 {
				server = bodyString[startIndex : startIndex+endIndex]
				util.Logger.Info("Auth server found", zap.String("server", server))
			}
		}

	default:
		util.Logger.Error("Find auth server: unexpected response", zap.String("state", response.Status))
	}
	return
}
