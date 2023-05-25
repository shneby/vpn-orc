package network

// todo: I'll need to refactor this

import (
	"fmt"
	"net"
)

type AddressPool struct {
	cidr    *net.IPNet
	used    map[string]bool
	network *net.IPNet
}

func NewAddressPool(cidr string) (*AddressPool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	used := make(map[string]bool)
	return &AddressPool{
		cidr:    ipNet,
		used:    used,
		network: ipNet,
	}, nil
}

func (ap *AddressPool) AllocateAddress() (net.IP, error) {
	for ip := ap.network.IP.Mask(ap.network.Mask); ap.network.Contains(ip); incrementIP(ip) {
		ipString := ip.String()
		if !ap.used[ipString] {
			ap.used[ipString] = true
			return ip, nil
		}
	}

	return nil, fmt.Errorf("no available addresses in the pool")
}

func (ap *AddressPool) ReturnAddress(ip net.IP) {
	ipString := ip.String()
	delete(ap.used, ipString)
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
