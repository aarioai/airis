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

// SQL can display it with select *, cast(ip as CHAR) from table_name
func ToIP(addr string) IP {
	var ip IP
	if addr == "" {
		return ip
	}
	nip := net.ParseIP(addr) // both IPv4 and IPv6 parse to 16 bytes
	if nip == nil {
		return ip
	}
	ip2 := nip.To4() // convert IPv4 to 4 bytes
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

// Whether this is IPv4
func (ip IP) Is4() bool {
	return len(ip.String) == net.IPv4len
}

// 4-byte IPv4 and 16-byte IPv4 or IPv6 all produce a readable IP address
func (ip IP) To16() string {
	ip2 := ip.Net()
	// includes IPv4 / IPv6
	nip := ip2.To16() // IP length is 16 here
	if nip == nil {
		return ""
	}
	return nip.String() // return an IPv4-style IP address
}
