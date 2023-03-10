package main

import (
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
	"github.com/cyp0633/drcom-go/internal/dhcp"
	dhcpauto "github.com/cyp0633/drcom-go/internal/dhcp/auto"
	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
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
			daemonize()
		}
		dhcp.Run()
	case "pppoe":
		util.Logger.Fatal("PPPoE mode not implemented")
	}
}

// 使程序在后台运行
func daemonize() {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		// 带 auto 的命令均用于生成配置文件，生成完毕后直接进行普通 DHCP 即可
		if args[i] == "dhcp-auto" {
			args[i] = "dhcp"
		}
		// 不要重复 fork 了
		if args[i] == "-d" || args[i] == "--daemon" {
			args = append(args[:i], args[i+1:]...)
		}
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Start()
	util.Logger.Info("Daemon started", zap.Int("pid", cmd.Process.Pid))
	os.Exit(0)
}
