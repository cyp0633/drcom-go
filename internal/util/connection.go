package util

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// CheckConnection 检查网络连接，如发现不通，向 ch 中发送信号
func CheckConnection(ch chan bool, interval time.Duration) {
	// 自定义 DNS
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}
			return dialer.DialContext(ctx, "udp", ExtConf.DnsServer+":53")
		},
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				ip, err := resolver.LookupHost(ctx, host)
				if err != nil {
					return nil, err
				}
				return net.Dial(network, net.JoinHostPort(ip[0], port))
			},
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	for {
		resp, err := client.Get(ExtConf.ConnectionTestServer)
		if err != nil || resp.StatusCode != http.StatusNoContent {
			ch <- true
			Logger.Warn("Network connection lost", zap.Error(err))
			return
		} else {
			resp.Body.Close()
			Logger.Debug("Network connection is OK")
			time.Sleep(interval)
		}
	}
}
