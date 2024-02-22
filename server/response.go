package server

import (
	"net/netip"
	"strconv"
)

type Response struct {
	IpAddress  string `json:"ip_addr"`
	RemoteHost string `json:"remote_host"`
	UserAgent  string `json:"user_agent"`
	Port       string `json:"port"`
	Language   string `json:"language"`
	Method     string `json:"method"`
	Encoding   string `json:"encoding"`
	Mime       string `json:"mime"`
	Via        string `json:"via"`
	Fowarded   string `json:"forwarded"`
}

func (r *Response) egress() string {
	ip, err := netip.ParseAddr(r.IpAddress)
	if err != nil {
		return ""
	}

	if port, err := strconv.Atoi(r.Port); err != nil {
		return ip.String()
	} else {
		return netip.AddrPortFrom(ip, uint16(port)).String()
	}
}
