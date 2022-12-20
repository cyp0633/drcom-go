package internal

var CLI struct {
	Dhcp    struct{} `cmd:"" help:"Use DHCP mode"`                                                    // DHCP模式
	Pppoe   struct{} `cmd:"" help:"Use PPPoE mode"`                                                   // PPPoE模式
	Conf    string   `help:"Configuration file path" short:"c" default:"/etc/drcom.conf" type:"path"` // 配置文件目录
	BindIP  string   `help:"IP address to bind to" short:"b" default:"0.0.0.0"`                       // 绑定IP地址
	Log     string   `help:"Log ONLY to specified path" short:"l" default:""`                         // 日志文件目录
	Daemon  bool     `help:"Run as daemon" short:"d" default:"false"`                                 // 是否在后台运行
	Eternal bool     `help:"Keep trying to reconnect" short:"e" default:"false"`                      // 是否一直尝试重连
}
