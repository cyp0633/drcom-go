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
	KeepAliveVersion   [2]byte
	RorVersion         bool
	Keepalive1Mod      bool
	PppoeFlag          byte
	MacParsed          [6]byte
	Keepalive2Flag     byte
}

// parseConf parse configuration file into Conf struct
func parseConf(iniPath string) {
	Conf.Keepalive1Mod = false // default value
	var temp string
	cfg, err := ini.Load(iniPath)
	if err != nil {
		log.Fatalln("Failed to read config file", err.Error())
	}
	sec := cfg.Section("")

	// universal
	Conf.Server = sec.Key("server").String()

	// DHCP only
	Conf.Username = sec.Key("username").String()
	Conf.Password = sec.Key("password").String()
	temp = sec.Key("CONTROLCHECKSTATUS").String()
	fmt.Sscanf(temp, "\\x%x", &Conf.ControlCheckStatus)
	temp = sec.Key("ADAPTERNUM").String()
	fmt.Sscanf(temp, "\\x%x", &Conf.AdapterNum)
	Conf.HostIp = sec.Key("host_ip").String()
	temp = sec.Key("IPDOG").String()
	fmt.Sscanf(temp, "\\x%x", &Conf.IpDog)
	Conf.Hostname = sec.Key("host_name").String()
	Conf.PrimaryDns = sec.Key("PRIMARY_DNS").String()
	Conf.DhcpServer = sec.Key("dhcp_server").String()
	temp = sec.Key("AUTH_VERSION").String()
	fmt.Sscanf(temp, "\\x%x\\x%x", &Conf.AuthVersion[0], &Conf.AuthVersion[1])
	Conf.Mac = sec.Key("mac").String()
	Conf.MacParsed = parseMac(sec.Key("mac").String())
	Conf.HostOs = sec.Key("host_os").String()
	temp = sec.Key("KEEP_ALIVE_VERSION").String()
	fmt.Sscanf(temp, "\\x%x\\x%x", &Conf.KeepAliveVersion[0], &Conf.KeepAliveVersion[1])
	Conf.RorVersion, _ = sec.Key("ror_version").Bool()
	Conf.Keepalive1Mod, _ = sec.Key("keepalive1_mod").Bool()

	// PPPoE only
	temp = sec.Key("pppoe_flag").String()
	fmt.Sscanf(temp, "\\x%x", Conf.PppoeFlag)
	temp = sec.Key("keep_alive2_flag").String()
	fmt.Sscanf(temp, "\\x%x", Conf.Keepalive2Flag)
}

// parseMac parse a mac address like xx:xx:xx:xx:xx:xx into byte array
func parseMac(mac string) [6]byte {
	var temp [6]byte
	var slice string
	for i := 0; i < 6; i++ {
		slice = mac[i*2+2 : i*2+4]
		_, err := fmt.Sscanf(slice, "%x", &temp[i])
		if err != nil {
			log.Panicln("Mac parse failed", err.Error())
		}
	}
	return temp
}
