package dhcp

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	"go.uber.org/zap"

	"github.com/cyp0633/drcom-go/internal/util"
)

// 第一个保活包
func keepAlive1(salt []byte, authInfo []byte) {
	// keepalive1_mod in dogcom?
	pkt := []byte{}
	pkt = append(pkt, 0x03, 0x01)                                                       // 0xff in drcom-generic?
	md5a := md5.Sum(append(append([]byte("\x03\x01"), salt...), util.Conf.Password...)) // salt = seed
	pkt = append(pkt, md5a[:]...)
	pkt = append(pkt, 0x00, 0x00, 0x00)
	pkt = append(pkt, authInfo...)
	pkt = append(pkt, byte(rand.Intn(0xFF)), byte(rand.Intn(0xFF)))
	_, err := conn.Write(pkt)
	if err == nil {
		err = conn.Flush()
	}
	if err != nil {
		util.Logger.Error("Sending keepalive1 packet failed", zap.Error(err))
		return
	}
	util.Logger.Debug("Keepalive1 sent", zap.String("packet", hex.EncodeToString(pkt)))

	// 读取keepalive1结果
	result := make([]byte, 1024)
	n, err := conn.Read(result)
	if err != nil {
		util.Logger.Error("Receiving keepalive1 result failed", zap.Error(err))
		return
	}
	util.Logger.Debug("Keepalive1 recv", zap.String("packet", hex.EncodeToString(result[:n])))
	if result[0] != 0x07 {
		util.Logger.Warn("Bad keepalive1 packet received", zap.String("packet", hex.EncodeToString(result[:n])))
	}
}

func keepAlive2(salt []byte, tail []byte) {

}
