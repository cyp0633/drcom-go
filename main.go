package main

import (
	"github.com/alecthomas/kong"
	"github.com/cyp0633/drcom-go/internal/dhcp"
	dhcpauto "github.com/cyp0633/drcom-go/internal/dhcp/auto"
	"github.com/cyp0633/drcom-go/internal/util"
)

func main() {
	ctx := kong.Parse(&util.CLI,
		kong.Description("Drcom client in Go, see https://github.com/cyp0633/drcom-go"),
		kong.UsageOnError(),
	)
	util.SetupLog()
	switch ctx.Command() {
	case "dhcp-auto":
		dhcpauto.Auto()
		fallthrough
	case "dhcp":
		util.ParseConf()
		if util.CLI.Daemon {
			util.Logger.Info("Daemon mode")
			util.Daemonize()
		}
		dhcp.Run()
	case "pppoe":
		util.Logger.Fatal("PPPoE mode not implemented")
	}
}
