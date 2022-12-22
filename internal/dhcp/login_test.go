package dhcp

import (
	"testing"

	"github.com/cyp0633/drcom-go/internal/util"
)

func TestGenLoginPacket(t *testing.T) {
	fillDummyData()
	util.SetupLog()
	p, err := genLoginPacket([]byte{0x98, 0x43, 0x75, 0x00})
	if err != nil {
		t.Error(err)
	}
	t.Log(p)
}

func fillDummyData() {
	util.Conf.Username = "13888888888@fku"
	util.Conf.Password = "123456"
	util.Conf.ControlCheckStatus = '\x20'
	util.Conf.AdapterNum = '\x05'
	util.Conf.Mac = "0x112233445566"
	util.Conf.MacBytes = []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	util.Conf.HostIP = "255.255.255.255"
	util.Conf.IpDog = '\x01'
	util.Conf.Hostname = "DESKTOP-123456"
	util.Conf.PrimaryDns = "8.8.8.8"
	util.Conf.DhcpServer = "10.0.0.25"
	util.Conf.AuthVersion = [2]byte{0x2a, 0x00}
	// incomplete
}