package dhcp

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	"go.uber.org/zap"

	"github.com/cyp0633/drcom-go/internal/util"
)

// 第一个保活包
func keepAlive1(salt []byte, authInfo []byte) error {
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
		return ErrorKeepalive1
	}
	util.Logger.Debug("Keepalive1 sent", zap.String("packet", hex.EncodeToString(pkt)))

	// 读取keepalive1结果
	result := make([]byte, 1024)
	n, err := conn.Read(result)
	if err != nil {
		util.Logger.Error("Receiving keepalive1 result failed", zap.Error(err))
		return ErrorKeepalive1
	}
	util.Logger.Debug("Keepalive1 recv", zap.String("packet", hex.EncodeToString(result[:n])))
	if result[0] != 0x07 {
		util.Logger.Warn("Bad keepalive1 packet received", zap.String("packet", hex.EncodeToString(result[:n])))
	}
	return nil
}

var keepAlive2Counter = 0

var first *int // to be improved

func keepAlive2(salt []byte, tail []byte) error {
	// send file packet
	if *first != 0 {
		pkt, err := genKeepalive2Packet(first, 1, 0)
		if err != nil {
			return ErrorKeepalive2
		}
		keepAlive2Counter++
		_, err = conn.Write(pkt)
		if err == nil {
			err = conn.Flush()
		}
		if err != nil {
			util.Logger.Error("Sending keepalive2 packet failed", zap.Error(err))
			return ErrorKeepalive2
		}
		util.Logger.Debug("Keepalive2_file sent", zap.String("packet", hex.EncodeToString(pkt)))

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			util.Logger.Error("Receiving keepalive2_file result failed", zap.Error(err))
			return ErrorKeepalive2
		}
		util.Logger.Debug("Keepalive2_file recv", zap.String("packet", hex.EncodeToString(buf[:n])))
		if buf[0] == 0x07 {
			if buf[2] == 0x10 {
				util.Logger.Debug("Authentic keepalive2_file recv")
			} else if buf[2] != 0x28 {
				util.Logger.Warn("Bad keepalive2_file packet received", zap.String("packet", hex.EncodeToString(buf[:n])))
				return ErrorKeepalive2
			}
		} else {
			util.Logger.Error("Bad keepalive2_file packet received", zap.String("packet", hex.EncodeToString(buf[:n])))
			return ErrorKeepalive2
		}
	}
	return nil
}

func genKeepalive2Packet(filepacket *int, typ, encryptType int) (pkt []byte, err error) {
	return
}
