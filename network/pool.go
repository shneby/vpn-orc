package network

import (
	"log"
	"net"
)

type Pool struct {
	cidr        *net.IPNet
	availableIP []net.IP
	allocatedIP []net.IP
}

func NewPool(cidr string) (*Pool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatalf("Unable to parse cidr [%s] - %s", cidr, err)
		return nil, err
	}

	pool := &Pool{
		cidr: ipNet,
	}

	// Generate the available IP addresses within the CIDR range

	ip := ipNet.IP
	log.Println(ip.IsGlobalUnicast())

	return pool, nil
}
