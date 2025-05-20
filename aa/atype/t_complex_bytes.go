package atype

import (
	"net"
)

func (p Position) Bytes() []byte {
	if !p.Ok() {
		return nil
	}
	return []byte(p.String)
}
func (p Position) Ok() bool {
	return p.Valid && len(p.String) == 25
}

// sql 可以使用 select *， cast(ip as CHAR) from table_name   进行显示
func ToIP(addr string) IP {
	var ip IP
	if addr == "" {
		return ip
	}
	nip := net.ParseIP(addr) // 无论是IPv4还是IPv6都是16字节
	if nip == nil {
		return ip
	}
	ip2 := nip.To4() // 将IPv4的转为4字节
	if ip2 != nil {
		nip = ip2
	}

	ip.Scan(nip.String())
	return ip
}
func (ip IP) Bytes() []byte {
	if !ip.OK() {
		return nil
	}
	return []byte(ip.String)
}

func (ip IP) OK() bool {
	return ip.Valid && len(ip.String) == net.IPv4len || len(ip.String) == net.IPv6len
}
func (ip IP) Net() net.IP {
	b := ip.Bytes()
	if b == nil {
		return nil
	}
	return b
}

// 是不是 IPv4
func (ip IP) Is4() bool {
	return len(ip.String) == net.IPv4len
}

// 无论是4字节的IPv4，还是16字节的IPv4或IPv6，都能输出可阅读的IP地址
func (ip IP) To16() string {
	ip2 := ip.Net()
	// 包括ipv4 / ipv16
	nip := ip2.To16() // 此时IP长度为16
	if nip == nil {
		return ""
	}
	return nip.String() // 返回IPv4样式IP地址
}
