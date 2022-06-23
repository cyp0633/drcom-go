package main

import (
	"github.com/cyp0633/go-drcom/network"
	util "github.com/cyp0633/go-drcom/util"
	"github.com/sevlyar/go-daemon"
	"log"
)

func main() {
	util.Parse()
	daemonCtx := &daemon.Context{
		PidFileName: "go-drcom.pid",
		PidFilePerm: 0644,
		LogFileName: util.Opt.LogPath,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{""},
	}
	if util.Opt.Daemon {
		d, err := daemonCtx.Reborn()
		if err != nil {
			log.Fatalln("Unable to run", err)
		}
		if d != nil {
			return
		}
		defer daemonCtx.Release()
	}
	if util.Opt.EnableX {
		network.TrySmartEaplogin()
	}
	network.Drcom(5)
}
