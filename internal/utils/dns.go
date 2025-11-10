package utils

import (
	"net"
)

// LookupIPv4 resolves the given domain and returns a slice of IPv4 addresses as strings.
func LookupIPv4(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	var ipv4s []string
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4s = append(ipv4s, ip.String())
		}
	}
	return ipv4s, nil
}
