package dhcpauto

import (
	"encoding/hex"
	"net"

	"github.com/cyp0633/drcom-go/internal/util"
	"go.uber.org/zap"
)

var serverList = []string{"119.39.30.2", "119.39.30.10", "119.39.30.18", "119.39.30.26", "119.39.30.34", "119.39.30.42", "119.39.32.42", "119.39.32.10", "119.39.32.66", "119.39.32.92", "119.39.32.110", "58.20.50.197", "58.20.127.10", "58.20.127.26", "58.20.127.34", "58.20.127.90", "58.20.127.100", "58.20.127.106", "58.20.127.178", "58.20.127.242", "58.20.127.250", "218.104.155.74", "10.255.255.2", "10.255.255.11", "119.39.20.2", "119.39.32.58", "119.39.32.154", "119.39.32.174", "119.39.32.202", "119.39.32.210", "172.17.0.253", "119.39.20.42", "58.20.41.227", "119.39.32.18", "119.39.32.82", "58.20.26.246", "119.39.21.18", "119.39.119.2", "119.39.119.66"}

var conn *net.UDPConn

// 向所有服务器发送探测包，谁回复了就是谁
func sendProbe() {
	addr := net.UDPAddr{IP: net.ParseIP(util.CLI.BindIP), Port: 61440}
	var err error
	conn, err = net.ListenUDP("udp", &addr)
	if err != nil {
		util.Logger.Fatal("Open socket on 61440 failed", zap.Error(err))
	}
	util.Logger.Debug("Opened socket on 61440")
	for _, server := range serverList {
		raddr := net.UDPAddr{IP: net.ParseIP(server), Port: 61440}
		_, err := conn.WriteToUDP([]byte("\x01\x02\xec\x97\x2a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), &raddr)
		if err != nil {
			util.Logger.Error("Send pkg to server failed", zap.Error(err))
		}
	}
}

// 接收探测包的回复
func recvProbe() {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, raddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		util.Logger.Error("Read from socket failed", zap.Error(err))
	}
	util.Logger.Debug("Read from socket", zap.String("data", hex.EncodeToString(buf[:n])))
	// check if raddr is in serverList
	for _, server := range serverList {
		if raddr.IP.String() == server {
			util.Conf.Server = server
			util.Logger.Debug("Found server", zap.String("server", server))
			return
		}
	}
	util.Logger.Error("Server not found", zap.String("fake_server", raddr.IP.String()))
}
