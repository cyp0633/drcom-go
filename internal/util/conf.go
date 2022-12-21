package util

import (
	"encoding/hex"

	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

// 对应 drcom-generic 配置
var Conf struct {
	Server             string
	Username           string
	Password           string
	ControlCheckStatus byte
	AdapterNum         byte
	HostIP             string
	IpDog              byte
	Hostname           string
	PrimaryDns         string
	DhcpServer         string
	AuthVersion        [2]byte
	Mac                string
	MacBytes           []byte
	HostOs             string
	KeepAliveVersion   [2]byte
	RorVersion         bool
}

// 解析配置文件
func ParseConf() {
	cfg, err := ini.Load(CLI.Conf)
	if err != nil {
		Logger.Panic("Opening configuration failed", zap.Error(err), zap.String("path", CLI.Conf))
	}
	var temp string
	section := cfg.Section("")
	Conf.Server = section.Key("server").String()
	Conf.Username = section.Key("username").String()
	Conf.Password = section.Key("password").String()
	temp = section.Key("CONTROLCHECKSTATUS").String()
	Conf.ControlCheckStatus = parseBytes(temp)[0] // 带有转义字符的字符串转换为 byte
	temp = section.Key("ADAPTERNUM").String()
	Conf.AdapterNum = parseBytes(temp)[0]
	Conf.HostIP = section.Key("host_ip").String()
	temp = section.Key("IPDOG").String()
	Conf.IpDog = parseBytes(temp)[0]
	Conf.Hostname = section.Key("host_name").String()
	Conf.PrimaryDns = section.Key("PRIMARY_DNS").String()
	Conf.DhcpServer = section.Key("dhcp_server").String()
	temp = section.Key("AUTH_VERSION").String()
	Conf.AuthVersion = [2]byte{parseBytes(temp)[0], parseBytes(temp)[1]}
	Conf.Mac = section.Key("mac").String()
	Conf.MacBytes, err = hex.DecodeString(Conf.Mac)
	if err != nil {
		Logger.Panic("Parsing conf mac failed", zap.Error(err), zap.String("mac", Conf.Mac))
	}
	Conf.HostOs = section.Key("host_os").String()
	temp = section.Key("KEEP_ALIVE_VERSION").String()
	Conf.KeepAliveVersion = [2]byte{parseBytes(temp)[0], parseBytes(temp)[1]}
	Conf.RorVersion = section.Key("ror_version").MustBool()
	Logger.Info("Configuration loaded", zap.String("path", CLI.Conf), zap.Any("conf", Conf))
}

// 带有转义字符的字符串转换为 byte slice
func parseBytes(s string) []byte {
	var b []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			switch s[i] {
			case 'x':
				h, err := hex.DecodeString(s[i+1 : i+3])
				if err != nil {
					Logger.Panic("Configuration bytes parsing failed", zap.Error(err), zap.String("Error byte", s[i:i+3]))
				}
				b = append(b, h[0])
				i += 2
			case 'r':
				b = append(b, '\r')
			case 'n':
				b = append(b, '\n')
			case 't':
				b = append(b, '\t')
			default:
				b = append(b, s[i])
			}
		} else {
			b = append(b, s[i])
		}
	}
	return b
}
