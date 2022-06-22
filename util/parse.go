package util

import (
	"fmt"
	"github.com/alecthomas/kingpin"
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

//Parse parse CLI arguments
func Parse() {
	//kingpin.Arg("mode", "set your drcom mode").Required().StringVar(&Opt.Mode)
	//kingpin.Arg("conf", "import configuration file").Required().StringVar(&Opt.Conf)
	//kingpin.Arg("bindip", "bind your ip address").Default("0.0.0.0").StringVar(&Opt.BindIp)
	//kingpin.Arg("log", "specify log file").StringVar(&Opt.LogPath)
	kingpin.Flag("mode", "set your drcom mode").Required().Short('m').StringVar(&Opt.Mode)
	kingpin.Flag("conf", "import configuration file").Required().Short('c').StringVar(&Opt.Conf)
	kingpin.Flag("bindip", "bind your ip address").Default("0.0.0.0").Short('b').StringVar(&Opt.BindIp)
	kingpin.Flag("log", "specify log file").Short('l').StringVar(&Opt.LogPath)
	kingpin.Flag("802.1x", "enable 802.11x").Short('x').BoolVar(&Opt.EnableX)
	kingpin.Flag("daemon", "set daemon flag").Short('d').BoolVar(&Opt.Daemon)
	kingpin.Flag("eternal", "set eternal flag").Short('e').BoolVar(&Opt.Eternal)
	kingpin.Flag("verbose", "set verbose flag").Short('v').BoolVar(&Opt.Verbose)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	if Opt.Mode != "dhcp" && Opt.Mode != "pppoe" {
		fmt.Println("Unknown mode")
		os.Exit(1)
	}
	parseConf(Opt.Conf)
}
