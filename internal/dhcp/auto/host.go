package dhcpauto

import (
	"net"
	"os"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

// 获取主机信息
func getHostInfo() {
	var err error
	util.Conf.Hostname, err = os.Hostname()
	if err != nil {
		util.Logger.Warn("Get hostname failed, using default value", zap.Error(err))
		util.Conf.Hostname = "Drcom"
	}
	util.Logger.Debug("Hostname", zap.String("hostname", util.Conf.Hostname))

	interfaces, err := net.Interfaces()
	if err != nil {
		util.Logger.Fatal("Get interfaces failed", zap.Error(err))
	}
	util.Conf.AdapterNum = byte(len(interfaces))
	util.Logger.Debug("Adapter number", zap.Int("adapterNum", int(util.Conf.AdapterNum)))

	util.Conf.HostIP = getIPInUse().String()
	util.Logger.Debug("Host IP", zap.String("hostIP", util.Conf.HostIP))

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			util.Logger.Fatal("Get interface addresses failed", zap.Error(err))
		}
		for _, addr := range addrs {
			if addr.String() == util.Conf.HostIP {
				util.Conf.Mac = i.HardwareAddr.String()
				util.Conf.MacBytes = i.HardwareAddr
				break
			}
		}
	}
	util.Logger.Debug("MAC address", zap.String("mac", util.Conf.Mac))

	// 官方客户端只有 Windows，你还能大方承认 Linux 不成？
	util.Conf.HostOs = "Windows 10"
}

// 获取连接认证服务器使用的 IP 地址和 MAC 地址
func getIPInUse() (ip net.IP) {
	// 强制绑定某个 IP 的话，当然就是它了
	if util.CLI.BindIP != "" {
		ip = net.ParseIP(util.CLI.BindIP)
	} else {
		// 获取本机连接认证服务器的 IP 地址
		conn, err := net.Dial("udp", util.Conf.Server+":80")
		if err != nil {
			util.Logger.Fatal("Dialing auth server failed. This is probably not your problem", zap.Error(err))
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		ip = localAddr.IP
	}
	return
}
