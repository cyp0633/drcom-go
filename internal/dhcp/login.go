package dhcp

import (
	"math/rand"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

// 登录
func login() (tail []byte, salt []byte, err error) {
	salt, err = challenge()
	return
}

// 发送 challenge 并获得 salt
func challenge() (salt []byte, err error) {
	pkt := make([]byte, 20)
	pkt[0] = 0x01
	pkt[1] = 0x02
	pkt[2] = byte(rand.Intn(0xff))
	pkt[3] = byte(rand.Intn(0xff))
	pkt[4] = util.Conf.AuthVersion[0]
	_, err = conn.Write(pkt)
	if err != nil {
		util.Logger.Error("Sending challenge failed", zap.Error(err))
		return nil, err
	}
	conn.Flush()
	util.Logger.Info("Challenge sent", zap.Uint8s("content", pkt))
	err = ErrorLogin
	return
}
