package dhcp

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"net/netip"
	"strings"

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
	util.Logger.Info("Challenge recv", zap.String("packet", hex.EncodeToString(salt[:n])), zap.String("salt", hex.EncodeToString(salt[4:8])))
	salt = salt[4:8] // 前一部分只有 [4:8] 不同，看起来有用
	return
}

// 合成登录数据包
func genLoginPacket(salt []byte) (loginPacket []byte, err error) {
	// _tagLoginPacket
	{
		// _tagDrComHeader
		loginPacket = append(loginPacket, 0x03, 0x01, 0x00, byte(len(util.Conf.Username)+20))
		// PasswordMD5
		md51 := md5.Sum(append(append([]byte{0x03, 0x01}, salt...), []byte(util.Conf.Password)...)) // 好 tm 丑陋
		loginPacket = append(loginPacket, md51[:]...)
		// Account
		zeros := make([]byte, 36-len(util.Conf.Username))
		loginPacket = append(loginPacket, append([]byte(util.Conf.Username), zeros...)...) // 看起来用户名区域必须填充到 36B
		// ControlCheckStatus
		loginPacket = append(loginPacket, util.Conf.ControlCheckStatus)
		// AdapterNum
		loginPacket = append(loginPacket, util.Conf.AdapterNum)
		// MacAddrXORPasswordMD5
		xor := binary.BigEndian.Uint64(loginPacket[2:10]) ^ binary.BigEndian.Uint64(append([]byte("\x00\x00"), util.Conf.MacBytes...)) // 最后只取 6 位，这里 8 位也没关系，注意大端序
		loginPacket = append(loginPacket, hex.EncodeToString([]byte{byte(xor >> 40), byte(xor >> 32), byte(xor >> 24), byte(xor >> 16), byte(xor >> 8), byte(xor)})...)
		// PasswordMd5_2
		md52 := md5.Sum(append(append(append([]byte{1}, []byte(util.Conf.Password)...), salt...), []byte{0, 0, 0, 0}...)) // 更 tm 丑陋了
		loginPacket = append(loginPacket, md52[:]...)
		// HostIpNum
		loginPacket = append(loginPacket, 0x01)
		// HostIPList。 后三个可以全填 0，即 12 个 0 字节
		strings.Split(util.Conf.HostIP, ".")
		var ip netip.Addr
		ip, err = netip.ParseAddr(util.Conf.HostIP)
		if err != nil {
			util.Logger.Error("Parse HostIP configuration failed", zap.Error(err), zap.String("HostIP", util.Conf.HostIP))
			err = ErrorGenLoginPacket
			return
		}
		loginPacket = append(loginPacket, ip.AsSlice()...)
		loginPacket = append(loginPacket, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // 三个空 hostip
		// HalfMD5
		md53 := md5.Sum(append(loginPacket, 0x14, 0x00, 0x07, 0x0b))
		loginPacket = append(loginPacket, md53[0:8]...)
		// DogFlag
		loginPacket = append(loginPacket, util.Conf.IpDog)
		// unknown2
		loginPacket = append(loginPacket, 0x00, 0x00, 0x00, 0x00)
		// _tagHostInfo
		{
			// HostName 填充到 32B
			zeros := make([]byte, 32-len(util.Conf.Hostname))
			loginPacket = append(loginPacket, append([]byte(util.Conf.Hostname), zeros...)...)
			// DNSIP1
			var ip netip.Addr
			ip, err = netip.ParseAddr(util.Conf.PrimaryDns)
			if err != nil {
				util.Logger.Error("Parse PrimaryDns configuration failed", zap.Error(err), zap.String("PrimaryDns", util.Conf.PrimaryDns))
				err = ErrorGenLoginPacket
				return
			}
			loginPacket = append(loginPacket, ip.AsSlice()...)
			// DHCPServerIP
			ip, err = netip.ParseAddr(util.Conf.DhcpServer)
			if err != nil {
				util.Logger.Error("Parse DhcpServer configuration failed", zap.Error(err), zap.String("DhcpServer", util.Conf.DhcpServer))
				err = ErrorGenLoginPacket
				return
			}
			loginPacket = append(loginPacket, ip.AsSlice()...)
			// DNSIP2 填充四个 0
			loginPacket = append(loginPacket, 0x00, 0x00, 0x00, 0x00)
			// WINSIP1、WINSIP2 各填充四个 0
			loginPacket = append(loginPacket, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
			// _tagOSVersionInfo
			{
				// OSVersionInfoSize
				loginPacket = append(loginPacket, 0x94, 0x00, 0x00, 0x00)
				// MajorVersion，与 MinorVersion 应该共同表示了 Windows NT 版本，此处为 Windows 10
				loginPacket = append(loginPacket, 0x10, 0x00, 0x00, 0x00)
				// MinorVersion
				loginPacket = append(loginPacket, 0x00, 0x00, 0x00, 0x00)
				// BuildNumber 小端序的 Windows Build 编号 10240
				loginPacket = append(loginPacket, 0x00, 0x28, 0x00, 0x00)
				// PlatformID
				loginPacket = append(loginPacket, 0x02, 0x00, 0x00, 0x00)
				// Service Pack 40 个字节从 mchome/dogcom 复制，再填充 16 字节到 64 字节，但似乎 drcoms/drcom-generic 规定的是 128 字节？
				loginPacket = append(loginPacket, 0x33, 0x64, 0x63, 0x37, 0x39, 0x66, 0x35, 0x32, 0x31, 0x32, 0x65, 0x38, 0x31, 0x37, 0x30, 0x61, 0x63, 0x66, 0x61, 0x39, 0x65, 0x63, 0x39, 0x35, 0x66, 0x31, 0x64, 0x37, 0x34, 0x39, 0x31, 0x36, 0x35, 0x34, 0x32, 0x62, 0x65, 0x37, 0x62, 0x31)
				loginPacket = append(loginPacket, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
			}
		}
		// ClientVerInfoAndInternetMode
		loginPacket = append(loginPacket, util.Conf.AuthVersion[:]...)
		// DogVersion 空 2 字节
		loginPacket = append(loginPacket, 0x00, 0x00)
		if util.Conf.RorVersion { // LDAPAuth
			// Code
			loginPacket = append(loginPacket, 0x00)
			// PasswordLen
			loginPacket = append(loginPacket, byte(len(util.Conf.Password)))
			// Password ROR
			ror := util.Ror(md51[:], []byte(util.Conf.Password))
			loginPacket = append(loginPacket, ror...)
		}
		// DrcomAuthExtData
		{
			// Code
			loginPacket = append(loginPacket, 0x02)
			// Len
			loginPacket = append(loginPacket, 0x0c)
			// CRC
			crc := util.Checksum(append(append(loginPacket, []byte{0x01, 0x26, 0x07, 0x11, 0x00, 0x00}...), util.Conf.MacBytes...))
			loginPacket = append(loginPacket, crc[:]...)
			// Option
			loginPacket = append(loginPacket, 0x00, 0x00)
			// AdapterAddress
			loginPacket = append(loginPacket, util.Conf.MacBytes...)
		}
		// auto logout(1), broadcast mode(1), unknown(2)
		loginPacket = append(loginPacket, 0x00, 0x00, 0xe9, 0x13)
	}
	return
}
