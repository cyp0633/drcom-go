package dhcpauto

import (
	"testing"

	"github.com/cyp0633/drcom-go/internal/util"
)

func TestGetIPInUse(t *testing.T) {
	util.Conf.Server = "39.156.66.10"
	util.CLI.BindIP = ""
	ip:=getIPInUse()
	t.Log(ip)
}
