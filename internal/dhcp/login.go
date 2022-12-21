package dhcp

import (
	"bytes"
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

	// 合成登录数据包
	loginPacket := []byte{}
	// _tagLoginPacket
	{
		// _tagDrComHeader
		loginPacket = append(loginPacket, 0x03, 0x01, 0x00, byte(len(util.Conf.Username)+20))
		// PasswordMD5

	}
	// 发送
	_, err = conn.Write(loginPacket)
	if err != nil {
		util.Logger.Error("Sending login packet failed", zap.Error(err))
		err = ErrorLogin
		return
	}
	conn.Flush()
	util.Logger.Debug("Login packet sent", zap.String("packet", hex.EncodeToString(loginPacket)))

	// 读取登录结果
	result := make([]byte, 1024)
	n, err := conn.Read(result)
	if err != nil {
		util.Logger.Error("Reading login result failed", zap.Error(err))
		err = ErrorLogin
		return
	}
	util.Logger.Debug("Login result recv", zap.String("packet", hex.EncodeToString(result[:n])))
	if bytes.Equal(result[0:1], []byte{0x04}) { // 登录成功
		util.Logger.Info("Logged in")
		tail = result[23:39] // 同时也是 authinfo
		return
	} else { // 登录失败
		util.Logger.Error("Login failed")
		if result[0] == 0x05 { // 使用 mchome/dogcom 的错误类型判断
			switch result[4] {
			case 0x01:
				util.Logger.Info("CHECK_MAC", zap.String("tip", "MAC address in use by another user"))
			case 0x02:
				util.Logger.Info("SERVER_BUSY", zap.String("tip", "Wait for a while and try again"))
			case 0x03:
				util.Logger.Info("WRONG_PASS", zap.String("tip", "Check your password"))
			case 0x04:
				util.Logger.Info("NOT_ENOUGH", zap.String("tip", "Check your time/traffic balance"))
			case 0x05:
				util.Logger.Info("FREEZE_UP", zap.String("tip", "Account suspended"))
			case 0x07:
				util.Logger.Info("NOT_ON_THIS_IP", zap.String("tip", "IP address is restricted"))
			case 0x11:
				util.Logger.Info("NOT_ON_THIS_MAC", zap.String("tip", "MAC address is restricted"))
			case 0x20:
				util.Logger.Info("TOO_MUCH_IP", zap.String("tip", "Too many devices in use"))
			case 0x21:
				util.Logger.Info("UPDATE_CLIENT", zap.String("tip", "The target version is not supported"))
			case 0x22:
				util.Logger.Info("NOT_ON_THIS_IP_MAC", zap.String("tip", "IP/MAC address is restricted"))
			case 0x23:
				util.Logger.Info("MUST_USE_DHCP", zap.String("tip", "Turn on DHCP instead of static IP"))
			}
		}
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
