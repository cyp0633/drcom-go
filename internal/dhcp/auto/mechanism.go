package dhcpauto

import (
	"net"

	"github.com/cyp0633/drcom-go/internal/util"
)

// 没错，没研究透，所以直接进行一个 guess
func guess(ver int) {
	switch ver {
	case Drcom52D:
		{
			util.Conf.ControlCheckStatus = '\x20'
			util.Logger.Warn("Not fully researched, assuming ControlCheckStatus = 0x20")

			util.Conf.IpDog = '\x01'
			util.Logger.Warn("Not fully researched, assuming IpDog = 0x01")

			util.Conf.PrimaryDns = "58.20.127.170"
			util.Logger.Warn("Not fully researched, assuming PrimaryDns = '58.20.127.170'")

			hostIp := net.ParseIP(util.Conf.HostIP)
			hostIp[4] = 0x00
			util.Conf.DhcpServer = hostIp.String()
			util.Logger.Warn("Not fully researched, assuming DhcpServer = " + util.Conf.DhcpServer)

			util.Conf.AuthVersion = [2]byte{0x2a, 0x00}
			util.Logger.Warn("Not fully researched, assuming AuthVersion = 0x2a00")

			util.Conf.KeepAliveVersion = [2]byte{0xd8, 0x02}
			util.Logger.Warn("Not fully researched, assuming KeepAliveVersion = 0xd802")

			util.Conf.RorVersion = false
			util.Logger.Warn("Not fully researched, assuming RorVersion = false")
		}
	case Drcom60D:
		{
			util.Logger.Panic("Not supported!")
		}
	}
}
