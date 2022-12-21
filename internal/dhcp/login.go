package dhcp

import (
	"encoding/hex"
	"math/rand"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

// 登录
func login() (tail []byte, salt []byte, err error) {
	salt, err = challenge()
	if err != nil {
		err = ErrorLogin
		return
	}
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
		err = ErrorChallenge
		return
	}
	conn.Flush()
	util.Logger.Info("Challenge sent", zap.String("packet", hex.EncodeToString(pkt)))
	// 读取 salt
	salt = make([]byte, 1024)
	n, err := conn.Read(salt)
	if err != nil {
		util.Logger.Error("Reading challenge salt failed", zap.Error(err))
		err = ErrorChallenge
		return
	}
	util.Logger.Info("Challenge recv", zap.String("packet", hex.EncodeToString(salt[:n])))
	salt = salt[4:8] // 前一部分只有 [4:8] 不同，看起来有用
	return
}
