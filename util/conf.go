package util

import (
	"github.com/go-ini/ini"
	"log"
)

var Conf struct {
	Server           string
	Username         string
	Password         string
	ControlCheck     byte
	AdapterNum       byte
	HostIp           string
	IpDog            byte
	Hostname         string
	PrimaryDns       string
	DhcpServer       string
	AuthVersion      string
	Mac              string
	HostOs           string
	KeepAliveVersion string
	RorVersion       bool
}

func parseConf(iniPath string) {
	cfg, err := ini.Load(iniPath)
	if err != nil {
		log.Fatalln("Failed to read config file", err.Error())
	}
	sec := cfg.Section("")
	Conf.Server = sec.Key("server").String()
	Conf.Username = sec.Key("username").String()
	Conf.Password = sec.Key("password").String()
	Conf.ControlCheck = sec.Key("CONTROLCHECKSTATUS").String()[0]
	Conf.AdapterNum = sec.Key("ADAPTERNUM").String()[0]
	Conf.HostIp = sec.Key("host_ip").String()
	Conf.IpDog = sec.Key("host_ip").String()[0]
	Conf.Hostname = sec.Key("host_name").String()
	Conf.PrimaryDns = sec.Key("PRIMARY_DNS").String()
	Conf.DhcpServer = sec.Key("dhcp_server").String()
	Conf.AuthVersion = sec.Key("AUTH_VERSION").String()
	Conf.Mac = sec.Key("mac").String()
	Conf.HostOs = sec.Key("host_os").String()
	Conf.KeepAliveVersion = sec.Key("KEEP_ALIVE_VERSION").String()
	Conf.RorVersion, _ = sec.Key("ror_version").Bool()
}
