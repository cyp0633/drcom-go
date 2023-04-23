package dhcp

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"net/netip"
	"strings"
	"time"

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
	loginPacket, err := genLoginPacket(salt)
	if err != nil {
		err = ErrorLogin
		return
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

	// 读取登录结果，限时 1s
	result := make([]byte, 1024)
	udpConn.SetDeadline(time.Now().Add(time.Second))
	n, err := conn.Read(result)
	udpConn.SetDeadline(time.Time{})
	if err != nil {
		util.Logger.Error("Reading login result failed", zap.Error(err))
		conn.Flush()
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
		err = ErrorLogin
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
	util.Logger.Debug("Challenge sent", zap.String("packet", hex.EncodeToString(pkt)))
	// 读取 salt
	salt = make([]byte, 1024)
	n, err := conn.Read(salt)
	if err != nil {
		util.Logger.Error("Reading challenge salt failed", zap.Error(err))
		err = ErrorChallenge
		return
	}
	util.Logger.Debug("Challenge recv", zap.String("packet", hex.EncodeToString(salt[:n])), zap.String("salt", hex.EncodeToString(salt[4:8])))
	// 长度应为 40 字节
	if len(salt) != 40 {
		util.Logger.Error("Challenge reply length does not match")
		err = ErrorChallenge
		return
	}
	salt = salt[4:8] // 前一部分只有 [4:8] 不同，看起来有用
	return
}

// 合成登录数据报
func genLoginPacket(salt []byte) (loginPacket []byte, err error) {
	if util.Conf.RorVersion {
		loginPacket = make([]byte, 338)
	} else {
		loginPacket = make([]byte, 330)
	}
	// LoginPacket
	{
		// DrComHeader 0-3
		loginPacket[0] = 0x03
		loginPacket[1] = 0x01
		loginPacket[2] = 0x00
		loginPacket[3] = byte(len(util.Conf.Username) + 20)

		// PasswordMD5 4-19
		md51 := md5.Sum(append(append([]byte{0x03, 0x01}, salt...), []byte(util.Conf.Password)...)) // 好 tm 丑陋
		copy(loginPacket[4:], md51[:])

		// Username 20-55
		copy(loginPacket[20:], []byte(util.Conf.Username))

		// ControlCheckStatus 56
		loginPacket[56] = util.Conf.ControlCheckStatus

		// AdapterNum 57
		loginPacket[57] = util.Conf.AdapterNum

		// MacAddrXORPasswordMD5 58-63
		xor := binary.BigEndian.Uint64(loginPacket[2:10]) ^ binary.BigEndian.Uint64(append([]byte("\x00\x00"), util.Conf.MacBytes...)) // 最后只取 6 位，这里 8 位也没关系，注意大端序
		binary.BigEndian.PutUint64(loginPacket[58:], xor)

		// PasswordMD5_2 64-79
		md52 := md5.Sum(append(append(append([]byte{1}, []byte(util.Conf.Password)...), salt...), []byte{0, 0, 0, 0}...)) // 更 tm 丑陋了
		copy(loginPacket[64:], md52[:])

		// HostIpNum 80
		loginPacket[80] = 0x01

		// HostIpList 81-96
		// Fill in only 1 IP here
		strings.Split(util.Conf.HostIP, ".")
		var ip netip.Addr
		ip, err = netip.ParseAddr(util.Conf.HostIP)
		if err != nil {
			util.Logger.Error("Parse HostIP configuration failed", zap.Error(err), zap.String("HostIP", util.Conf.HostIP))
			err = ErrorGenLoginPacket
			return
		}
		copy(loginPacket[81:], ip.AsSlice())

		// HalfMD5 97-104
		md53 := md5.Sum(append(loginPacket[:97], 0x14, 0x00, 0x07, 0x0b))
		copy(loginPacket[97:], md53[:4])

		// DogFlag 105
		loginPacket[105] = util.Conf.IpDog

		// HostInfo
		{
			// Hostname 110-141
			copy(loginPacket[110:], []byte(util.Conf.Hostname))

			// PrimaryDNS 142-145
			var ip netip.Addr
			ip, err = netip.ParseAddr(util.Conf.PrimaryDns)
			if err != nil {
				util.Logger.Error("Parse PrimaryDns configuration failed", zap.Error(err), zap.String("PrimaryDns", util.Conf.PrimaryDns))
				err = ErrorGenLoginPacket
				return
			}
			copy(loginPacket[142:], ip.AsSlice())

			// DHCP 146-149
			ip, err = netip.ParseAddr(util.Conf.DhcpServer)
			if err != nil {
				util.Logger.Error("Parse DhcpServer configuration failed", zap.Error(err), zap.String("DhcpServer", util.Conf.DhcpServer))
				err = ErrorGenLoginPacket
				return
			}
			copy(loginPacket[146:], ip.AsSlice())

			// OSVersionInfo
			{
				// OSVersionInfoSize 162-165
				loginPacket[162] = 0x94
				loginPacket[163] = 0x00
				loginPacket[164] = 0x00
				loginPacket[165] = 0x00

				// MajorVersion 166-169
				loginPacket[166] = 0x10
				loginPacket[167] = 0x00
				loginPacket[168] = 0x00
				loginPacket[169] = 0x00

				// MinorVersion 170-173
				loginPacket[170] = 0x00
				loginPacket[171] = 0x00
				loginPacket[172] = 0x00
				loginPacket[173] = 0x00

				// BuildNumber 174-177
				loginPacket[174] = 0x00
				loginPacket[175] = 0x28
				loginPacket[176] = 0x00
				loginPacket[177] = 0x00

				// PlatformID 178-181
				loginPacket[178] = 0x02
				loginPacket[179] = 0x00
				loginPacket[180] = 0x00
				loginPacket[181] = 0x00
			}

		}
		// AuthVersion 310-311
		copy(loginPacket[310:], util.Conf.AuthVersion[:])

		// ExtData
		{
			// Code 312
			loginPacket[312] = 0x02

			// Len
			loginPacket[313] = 0x0c

			// CRC 314-317
			crc := util.Checksum(append(append(loginPacket, []byte{0x01, 0x26, 0x07, 0x11, 0x00, 0x00}...), util.Conf.MacBytes...))
			copy(loginPacket[314:], crc[:])

			// AdapterAddress 320-325
			copy(loginPacket[320:], util.Conf.MacBytes)
		}

		// BroadcastMode 326
		loginPacket[326] = 0xe9

		// Unknown 327
		loginPacket[327] = 0x13
	}

	return
}
