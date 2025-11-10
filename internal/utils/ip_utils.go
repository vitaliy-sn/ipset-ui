package utils

import (
	"net"
)

// Is the string a CIDR
func IsValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// Is the string an IP address
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// Does the IP address fall within the CIDR
func IsIPInCIDR(ip, cidr string) bool {
	ipAddr := net.ParseIP(ip)
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	return ipNet.Contains(ipAddr)
}

// Does one CIDR fall within another
func IsCIDRInCIDR(cidr1, cidr2 string) bool {
	_, ipNet1, err1 := net.ParseCIDR(cidr1)
	_, ipNet2, err2 := net.ParseCIDR(cidr2)
	if err1 != nil || err2 != nil {
		return false
	}
	return ipNet2.Contains(ipNet1.IP)
}
