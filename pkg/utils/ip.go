package utils

import (
	"fmt"
	"net"
	"time"
)

var (
	localPublicIP     string
	defaultUDPServers = []string{
		"223.5.5.5:80",       // 阿里云公共DNS
		"114.114.114.114:80", // 114公共DNS
		"8.8.8.8:80",         // Google公共DNS
	}
)

// LocalPublicIP 获取本机公网IP
func LocalPublicIP(udpServers []string, renew bool) (string, error) {
	if localPublicIP != "" && !renew {
		return localPublicIP, nil
	}

	if len(udpServers) == 0 {
		udpServers = defaultUDPServers
	}

	for _, server := range udpServers {
		ip, err := tryGetIP(server)
		if err == nil && ip != "" {
			localPublicIP = ip
			return localPublicIP, nil
		}
	}

	return "", fmt.Errorf("failed to get local public IP from all servers")
}

// tryGetIP 尝试从单个服务器获取IP
func tryGetIP(server string) (string, error) {
	conn, err := net.DialTimeout("udp", server, time.Millisecond*100)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return addr.IP.String(), nil
	}

	return "", fmt.Errorf("invalid address type")
}
