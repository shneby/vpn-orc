package network

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"net"
)

// note: The address pool pre-allocates the entire mask range - given more time we should implement sparse allocation

type AddressPool struct {
	cidr    *net.IPNet
	avail   stack.Stack
	used    map[string]bool
	network *net.IPNet
}

func NewAddressPool(cidr string) (*AddressPool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	used := make(map[string]bool)
	pool := &AddressPool{
		cidr:    ipNet,
		avail:   stack.Stack{},
		used:    used,
		network: ipNet,
	}
	pool.generateAddresses()
	return pool, nil
}

func (ap *AddressPool) generateAddresses() {
	ip := ap.network.IP.Mask(ap.network.Mask)

	ip = nextIP(ip) // skip the network address
	for ap.network.Contains(ip) {
		ap.avail.Push(ip)
		ip = nextIP(ip)
	}
	ap.avail.Pop() // pop the broadcast address from the stack
}

func nextIP(ip net.IP) net.IP {
	i := ip.To4()                                                     // break down the ip address to a 4 cell []byte
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3]) // convert to a uint and bit shift each cell to it's octet location
	v++                                                               // Increment the resulting number by one
	v3 := byte(v & 0xFF)                                              // break down each octet
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3) // recreate the address as an inet.IP
}

func (ap *AddressPool) AllocateAddress() (string, error) {
	if ap.avail.Len() == 0 {
		return "", fmt.Errorf("no available addresses in the pool")
	}
	ipString := fmt.Sprint(ap.avail.Pop())
	ap.used[ipString] = true
	return ipString, nil
}

func (ap *AddressPool) DeallocateAddress(address string) {
	if ap.used[address] {
		ap.avail.Push(net.ParseIP(address))
		delete(ap.used, address)
	}
}
