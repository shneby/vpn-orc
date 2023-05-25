package network

import (
	"errors"
	"log"
	"vpn-orc/persistence"
)

// The responsibilities of this services are as follows
// Onboarding a peer/tenant
// 1. The service assigns an address pool when a new tenant is introduced
// 2. when a peer of a tenant is introduced the tenants pool is used to assign an ip address from that pool
// 3. notifies an existing peers of the new peer
// Offboarding a peer
// 1. when a revoke is requested the service removes the peer, revokes the ip address and returns it to the network pool
// 2. notifies via notification service that the peer was revoked

type AddressService struct {
	// Repository service (represented by two maps right now)
	//addressPool     *DummyPool
	tenantToNetwork map[string]*DummyPool
	tenantToPeers   map[string][]persistence.Peer
	// Notification service
}

type AddressInterface interface {
	AllocateAddress(tenantId string, peerId string, publicKey string) (string, error)
	RevokeAddress(tenantId string, peerId string) error
	Main()
}

func NewAddressService() AddressInterface {
	return &AddressService{
		tenantToNetwork: make(map[string]*DummyPool),
		tenantToPeers:   make(map[string][]persistence.Peer),
	}
}

func (a *AddressService) AllocateAddress(tenantId string, peerId string, publicKey string) (string, error) { // todo: implement peer type

	if a.tenantToNetwork[tenantId] == nil {
		a.tenantToNetwork[tenantId] = NewDummyPool("192.168.1.0/24") // todo: this needs to be randomly generated somehow
	}

	tenantNetwork := a.tenantToNetwork[tenantId]
	address, err := tenantNetwork.AllocateIP()
	if err != nil {
		return "", err
	}

	log.Printf("Allocating address [%s] to peer [%s] over tenant [%s]", address, peerId, tenantId)
	return address, nil
}

func (a *AddressService) RevokeAddress(tenantId string, peerId string) error {

	if a.tenantToNetwork[tenantId] == nil {
		return errors.New("tenant " + tenantId + " has no allocated networks")
	}

	tenantNetwork := a.tenantToNetwork[tenantId]
	peerAddress := "192.168.1.1"

	err := tenantNetwork.DeallocateIP(peerAddress)
	if err != nil {
		return err
	}

	log.Println("Revoking address from connection pool")
	log.Println("Persist to database")
	log.Println("Read peers belonging tenantId from database")
	log.Println("Send notification to all peers on removed peer")

	return nil
}

func (a *AddressService) Main() {
	a.AllocateAddress("tenant1", "peer1", "adsada")
	a.AllocateAddress("tenant1", "peer1", "adsada")
	a.AllocateAddress("tenant2", "peer1", "adsada")

	err := a.RevokeAddress("tenant1", "peer1")
	if err != nil {
		log.Println(err)
	}
}
