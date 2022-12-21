package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
	"github.com/cyp0633/drcom-go/internal/dhcp"
	"github.com/cyp0633/drcom-go/internal/util"
)

func main() {
	ctx := kong.Parse(&util.CLI,
		kong.Description("Drcom client in Go, see https://github.com/cyp0633/drcom-go"),
		kong.UsageOnError(),
	)
	if util.CLI.Daemon {
		fmt.Printf("daemon mode\n")
		daemonize()
	}
	util.SetLogPath()
	util.ParseConf()
	if util.CLI.Eternal {
		fmt.Printf("eternal mode\n")
	}
	fmt.Printf("conf path: %s\n", util.CLI.Conf)
	fmt.Printf("bind ip: %s\n", util.CLI.BindIP)
	fmt.Printf("log path: %s\n", util.CLI.Log)
	switch ctx.Command() {
	case "dhcp":
		dhcp.Run()
	case "pppoe":
		util.Logger.Fatal("PPPoE mode not implemented")
	}
}

// 使程序在后台运行
func daemonize() {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "-d" || args[i] == "--daemon" {
			args[i] = ""
			break
		}
	}
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Start()
	fmt.Println("Daemon started with pid", cmd.Process.Pid)
	os.Exit(0)
}
