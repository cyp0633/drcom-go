package util

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
)

var Conf struct {
	Server             string
	Username           string
	Password           string
	ControlCheckStatus byte
	AdapterNum         byte
	HostIp             string
	IpDog              byte
	Hostname           string
	PrimaryDns         string
	DhcpServer         string
	AuthVersion        [2]byte
	Mac                string // unparsed mac address
	HostOs             string
	KeepAliveVersion   string
	RorVersion         bool
	Keepalive1Mod      int
	PppoeFlag          byte
	MacParsed          [6]byte
}

// parseConf parse configuration file into Conf struct
func parseConf(iniPath string) {
	Conf.Keepalive1Mod = 0 // default value
	var temp string
	cfg, err := ini.Load(iniPath)
	if err != nil {
		log.Fatalln("Failed to read config file", err.Error())
	}
	sec := cfg.Section("")
	Conf.Server = sec.Key("server").String()
	Conf.Username = sec.Key("username").String()
	Conf.Password = sec.Key("password").String()
	temp = sec.Key("CONTROLCHECKSTATUS").String()
	fmt.Sscanf(temp, "%02x", &Conf.ControlCheckStatus)
	temp = sec.Key("ADAPTERNUM").String()
	fmt.Sscanf(temp, "%02x", &Conf.AdapterNum)
	Conf.HostIp = sec.Key("host_ip").String()
	temp = sec.Key("IPDOG").String()
	fmt.Sscanf(temp, "%02x", &Conf.IpDog)
	Conf.Hostname = sec.Key("host_name").String()
	Conf.PrimaryDns = sec.Key("PRIMARY_DNS").String()
	Conf.DhcpServer = sec.Key("dhcp_server").String()
	temp = sec.Key("AUTH_VERSION").String()
	fmt.Sscanf(temp, "%02x%02x", &Conf.AuthVersion[0], &Conf.AuthVersion[1])
	fmt.Sscanf(temp, "%02x%02x", &Conf.AuthVersion[0], &Conf.AuthVersion[1])
	//Conf.Mac = sec.Key("mac").String()
	Conf.MacParsed = parseMac(sec.Key("mac").String())
	Conf.HostOs = sec.Key("host_os").String()
	Conf.KeepAliveVersion = sec.Key("KEEP_ALIVE_VERSION").String()
	Conf.RorVersion, _ = sec.Key("ror_version").Bool()
}

// parseMac parse a mac address like xx:xx:xx:xx:xx:xx into byte array
func parseMac(mac string) [6]byte {
	var temp [6]byte
	_, err := fmt.Sscanf(mac, "%02x%02x%02x%02x%02x%02x", &temp[0], &temp[1], &temp[2], &temp[3], &temp[4], &temp[5])
	if err != nil {
		log.Panicln("Mac parse failed", err.Error())
	}
	return temp
}
