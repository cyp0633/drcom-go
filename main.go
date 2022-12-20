package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
	cli "github.com/cyp0633/drcom-go/internal/cli"
	"github.com/cyp0633/drcom-go/internal/util"
)

func main() {
	ctx := kong.Parse(&cli.CLI,
		kong.Description("Drcom client in Go, see https://github.com/cyp0633/drcom-go"),
		kong.UsageOnError(),
	)
	if cli.CLI.Daemon {
		fmt.Printf("daemon mode\n")
		daemonize()
	}
	util.SetLogPath()
	if cli.CLI.Eternal {
		fmt.Printf("eternal mode\n")
	}
	fmt.Printf("conf path: %s\n", cli.CLI.Conf)
	fmt.Printf("bind ip: %s\n", cli.CLI.BindIP)
	fmt.Printf("log path: %s\n", cli.CLI.Log)
	switch ctx.Command() {
	case "dhcp":
		fmt.Printf("dhcp mode\n")
	case "pppoe":
		fmt.Printf("pppoe mode\n")
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
