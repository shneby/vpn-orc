package network

import (
	"errors"
	"fmt"
	"log"
	"vpn-orc/persistence"
)

type AddressService struct {
	repo                  persistence.RepositoryInterface
	tenantIdToAddressPool map[string]*AddressPool
	// Notification service
}

type AddressInterface interface {
	AllocateAddress(tenantId string, peerId string, publicKey []byte) (*persistence.OnboardingResponse, error)
	RevokeAddress(tenantId string, peerId string) error
}

func NewAddressService(repo persistence.RepositoryInterface) AddressInterface {
	tenants, err := repo.ReadTenants()
	if err != nil {
		log.Fatalf("unable to retreive tenants from database - %s", err)
	}

	service := &AddressService{
		repo:                  repo,
		tenantIdToAddressPool: make(map[string]*AddressPool),
	}

	for _, tenant := range tenants {
		addressPool, _ := NewAddressPool(tenant.Network)
		if err != nil {
			log.Fatalf("unable to instantiate network pool for tenant [%s]", tenant.Id)
		}
		// todo: query for all tenant peers and modify addressPool.used to reflect that the addresses are already in use.
		service.tenantIdToAddressPool[tenant.Id] = addressPool
	}

	return service
}

// AllocateAddress todo: flow of this method is big - need to move some logic towards the controller
func (a *AddressService) AllocateAddress(tenantId string, peerId string, publicKey []byte) (*persistence.OnboardingResponse, error) {
	log.Printf("received allocation request: tenant [%s], peer [%s]", tenantId, peerId)

	// Check tenant exist
	_, err := a.repo.ReadTenant(tenantId)
	if err != nil {
		return nil, errors.New("Unable to find tenant " + tenantId)
	}

	// Check peer exists on tenant
	peer, err := a.repo.ReadPeer(tenantId, peerId)
	if peer != nil {
		return nil, fmt.Errorf("peer %s already has an address", peer.Id)
	}

	// Create new peer object and write to db
	networkPool := a.tenantIdToAddressPool[tenantId]
	address, err := networkPool.AllocateAddress()
	if err != nil {
		return nil, fmt.Errorf("unable to allocate addresss for peer %s - %s", peerId, err)
	}

	// get all tenant peers before adding the new peer
	peers, _ := a.repo.ReadPeers(tenantId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch peers for tenant %s - %s", tenantId, err)
	}

	peerDTO := persistence.Peer{Id: peerId, Address: address, PublicKey: publicKey}
	err = a.repo.WritePeer(tenantId, peerDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to persist peer %s - %s", peer.Id, err)
	}

	// Add peer entry to health service
	// todo: implement this

	// send notification peer joined
	// todo: implement this

	// return address response with allocated address to peer
	return &persistence.OnboardingResponse{
		Address: address,
		Peers:   peers,
	}, nil
}

func (a *AddressService) RevokeAddress(tenantId string, peerId string) error {
	tenant, err := a.repo.ReadTenant(tenantId)
	if err != nil {
		return errors.New("Unable to find network for tenant " + tenantId)
	}

	log.Printf("RevokeAddress: tenantId[%s], peer[%s], network[%s]", tenantId, peerId, tenant.Network)
	return nil
}
