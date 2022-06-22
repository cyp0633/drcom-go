package util

import (
	"flag"
	"fmt"
	"os"
)

var Opt struct {
	Mode    string
	Conf    string
	BindIp  string
	LogPath string
	EnableX bool
	Daemon  bool
	Eternal bool
	Verbose bool
	Help    bool
}

func Parse() {
	flag.StringVar(&Opt.Mode, "m", "", "set your drcom Mode")
	if Opt.Mode != "dhcp" && Opt.Mode != "pppoe" {
		fmt.Println("Unknown mode")
		os.Exit(1)
	}
	flag.StringVar(&Opt.Conf, "c", "", "import configuration file")
	parseConf(Opt.Conf)
	flag.StringVar(&Opt.BindIp, "b", "0.0.0.0", "bind your ip address")
	flag.StringVar(&Opt.LogPath, "l", "", "specify log file")
	flag.BoolVar(&Opt.EnableX, "x", false, "enable 802.11x")
	flag.BoolVar(&Opt.Daemon, "d", false, "set daemon flag (Unix-like only)")
	flag.BoolVar(&Opt.Eternal, "e", false, "set eternal flag")
	flag.BoolVar(&Opt.Verbose, "v", false, "set verbose flag")
	flag.BoolVar(&Opt.Help, "h", false, "show help")
	flag.Parse()
	if Opt.Help {
		flag.Usage()
		os.Exit(0)
	}
}
