package util

import (
	"os"
	"os/exec"

	"go.uber.org/zap"
)

var CLI struct {
	Dhcp     struct{} `cmd:"" help:"Use DHCP mode"`                                                    // DHCP模式
	DhcpAuto struct{} `cmd:"" help:"Use DHCP mode with auto configuration"`                            // DHCP自动配置模式
	Pppoe    struct{} `cmd:"" help:"Use PPPoE mode"`                                                   // PPPoE模式
	Conf     string   `help:"Configuration file path" short:"c" default:"/etc/drcom.conf" type:"path"` // 配置文件目录
	BindIP   string   `help:"IP address to bind to" short:"b" default:"0.0.0.0"`                       // 绑定IP地址
	Log      string   `help:"Log ONLY to specified path" short:"l" default:""`                         // 日志文件目录
	Daemon   bool     `help:"Run as daemon" short:"d" default:"false"`                                 // 是否在后台运行
	Eternal  bool     `help:"Keep trying to reconnect" short:"e" default:"false"`                      // 是否一直尝试重连
	Debug    bool     `help:"Print debug level logging" short:"D" default:"false"`                     // 是否打印调试信息
}

// 使程序在后台运行
func Daemonize() {
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
	Logger.Info("Daemon started", zap.Int("pid", cmd.Process.Pid))
	os.Exit(0)
}
