package util

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

// 对应 drcom-generic 配置
type BaseConf struct {
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

var Conf BaseConf

// drcom-go 专有扩展配置
var ExtConf struct {
	// 检查连接（204）使用的 URL
	ConnectionTestServer string
	// 检查连接使用的 DNS
	DnsServer string
	// 断连行为
	ActionOnDisconnect int
	// 连接成功后执行的外部命令
	ExecOnConnect string
	// 连接失败后执行的外部命令
	ExecOnDisconnect string
}

// 断连行为
const (
	// 重试连接
	DisconnectActionRetry = iota
	// 重启程序
	DisconnectActionRestart
)

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
	Conf.MacBytes, err = hex.DecodeString(Conf.Mac[2:])
	if err != nil {
		Logger.Panic("Parsing conf mac failed", zap.Error(err), zap.String("mac", Conf.Mac))
	}
	Conf.HostOs = section.Key("host_os").String()
	temp = section.Key("KEEP_ALIVE_VERSION").String()
	Conf.KeepAliveVersion = [2]byte{parseBytes(temp)[0], parseBytes(temp)[1]}
	Conf.RorVersion = section.Key("ror_version").MustBool()
	Logger.Debug("Configuration loaded", zap.String("path", CLI.Conf), zap.Any("conf", Conf))

	// 扩展配置
	ext := cfg.Section("extend")
	if temp = ext.Key("connection_test_server").String(); temp != "" {
		ExtConf.ConnectionTestServer = temp
	} else {
		ExtConf.ConnectionTestServer = "https://connect.rom.miui.com/generate_204"
	}
	if temp = ext.Key("dns_server").String(); temp != "" {
		ExtConf.DnsServer = temp
	} else {
		ExtConf.DnsServer = "8.8.8.8"
	}
	temp = ext.Key("action_on_disconnect").String()
	switch temp {
	case "retry":
		ExtConf.ActionOnDisconnect = DisconnectActionRetry
	case "restart":
		ExtConf.ActionOnDisconnect = DisconnectActionRestart
	default:
		ExtConf.ActionOnDisconnect = DisconnectActionRetry
	}
	ExtConf.ExecOnConnect = ext.Key("exec_on_connect").String()
	ExtConf.ExecOnDisconnect = ext.Key("exec_on_disconnect").String()
	Logger.Debug("Extended configuration loaded", zap.Any("conf", ExtConf))
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

// 生成配置文件
func (c *BaseConf) SaveConf(path string) {
	cfg := ini.Empty()
	section := cfg.Section("")
	section.Key("server").SetValue("'" + c.Server + "'")
	section.Key("username").SetValue("'" + c.Username + "'")
	section.Key("password").SetValue("'" + c.Password + "'")
	section.Key("CONTROLCHECKSTATUS").SetValue(fmt.Sprintf("'\\x%02x'", c.ControlCheckStatus))
	section.Key("ADAPTERNUM").SetValue(fmt.Sprintf("'\\x%02x'", c.AdapterNum))
	section.Key("host_ip").SetValue("'" + c.HostIP + "'")
	section.Key("IPDOG").SetValue(fmt.Sprintf("'\\x%02x'", c.IpDog))
	section.Key("host_name").SetValue("'" + c.Hostname + "'")
	section.Key("PRIMARY_DNS").SetValue("'" + c.PrimaryDns + "'")
	section.Key("dhcp_server").SetValue("'" + c.DhcpServer + "'")
	section.Key("AUTH_VERSION").SetValue(fmt.Sprintf("'\\x%02x\\x%02x'", c.AuthVersion[0], c.AuthVersion[1]))
	section.Key("mac").SetValue("'" + c.Mac + "'")
	section.Key("host_os").SetValue("'" + c.HostOs + "'")
	section.Key("KEEP_ALIVE_VERSION").SetValue(fmt.Sprintf("'\\x%02x\\x%02x'", c.KeepAliveVersion[0], c.KeepAliveVersion[1]))
	section.Key("ror_version").SetValue("'" + strconv.FormatBool(c.RorVersion) + "'")
	err := cfg.SaveTo(path)
	if err != nil {
		Logger.Panic("Saving configuration failed", zap.Error(err), zap.String("path", path))
	} else {
		Logger.Info("Configuration saved", zap.String("path", path))
	}
}
