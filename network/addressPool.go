package network

// todo: TASK 	- refactor this code to match a single coherent style
// todo: BUG 	- address pool assigns network & broadcast addresses

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

func (ap *AddressPool) AllocateAddress() (string, error) {
	for ip := ap.network.IP.Mask(ap.network.Mask); ap.network.Contains(ip); incrementIP(ip) {
		ipString := ip.String()
		if !ap.used[ipString] {
			ap.used[ipString] = true
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("no available addresses in the pool")
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
