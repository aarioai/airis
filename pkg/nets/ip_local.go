package nets

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

// LocalPublicIP public IP
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

// LanIP Local Area Network IP, A (10.x.x.x), B (172.16.x.x - 172.31.x.x), C (192.168.x.x)
func LanIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, inf := range interfaces {
		// Ignore loop back IP, 127.x.x.x
		if inf.Flags&net.FlagLoopback != 0 || inf.Flags&net.FlagUp == 0 {
			continue
		}
		// ignore Docker virtual network
		if inf.Name == "docker0" || inf.Name == "br-" {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if ok && ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no LAN IP found")
}
