package network

import "errors"

type DummyPool struct {
	cidr        string
	availableIP []string
	allocatedIP []string
}

func NewDummyPool(cidr string) *DummyPool {

	pool := &DummyPool{
		cidr:        cidr,
		availableIP: []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"},
		allocatedIP: []string{},
	}

	return pool
}

func (d *DummyPool) AllocateIP() (string, error) {
	if len(d.availableIP) > 0 {
		address := d.availableIP[len(d.availableIP)-1]
		d.availableIP = d.availableIP[:len(d.availableIP)-1]
		d.allocatedIP = append(d.allocatedIP, address)
		return address, nil
	}
	return "", errors.New("no available addresses in network " + d.cidr)
}

func (d *DummyPool) DeallocateIP(ip string) error {

	for i, address := range d.allocatedIP {
		if address == ip {
			d.allocatedIP[i] = d.allocatedIP[len(d.allocatedIP)-1]
			d.allocatedIP = d.allocatedIP[:len(d.allocatedIP)-1]
		}
	}

	return errors.New("Address " + ip + " not allocated in network " + d.cidr)
}
