package util

import (
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

// 对应 drcom-generic 配置
var Conf struct {
	Server             string
	Username           string
	Password           string
	ControlCheckStatus string
	AdapterNum         string
	HostIP             string
	IpDog              string
	Hostname           string
	PrimaryDns         string
	DhcpServer         string
	AuthVersion        string
	Mac                string
	HostOs             string
	KeepAliveVersion   string
	RorVersion         bool
}

// 解析配置文件
func ParseConf() {
	cfg, err := ini.Load(CLI.Conf)
	if err != nil {
		Logger.Panic("Opening configuration failed", zap.Error(err))
	}
	section := cfg.Section("")
	Conf.Server = section.Key("server").String()
	Conf.Username = section.Key("username").String()
	Conf.Password = section.Key("password").String()
	Conf.ControlCheckStatus = section.Key("CONTROLCHECKSTATUS").String()
	Conf.AdapterNum = section.Key("ADAPTERNUM").String()
	Conf.HostIP = section.Key("host_ip").String()
	Conf.IpDog = section.Key("IPDOG").String()
	Conf.Hostname = section.Key("host_name").String()
	Conf.PrimaryDns = section.Key("PRIMARY_DNS").String()
	Conf.DhcpServer = section.Key("dhcp_server").String()
	Conf.AuthVersion = section.Key("AUTH_VERSION").String()
	Conf.Mac = section.Key("mac").String()
	Conf.HostOs = section.Key("host_os").String()
	Conf.KeepAliveVersion = section.Key("KEEP_ALIVE_VERSION").String()
	Conf.RorVersion = section.Key("ror_version").MustBool()
	Logger.Info("Configuration loaded", zap.Any("conf", Conf))
}
