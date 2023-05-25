package network

import (
	"errors"
	"log"
	"vpn-orc/persistence"
)

type AddressService struct {
	repo                  persistence.RepositoryInterface
	tenantIdToAddressPool map[string]*AddressPool
	// Notification service
}

type AddressInterface interface {
	AllocateAddress(tenantId string, peerId string, publicKey string) error
	RevokeAddress(tenantId string, peerId string) error
}

func NewAddressService(repo persistence.RepositoryInterface) AddressInterface {
	tenants, err := repo.ReadTenants()
	if err != nil {
		log.Fatalf("Unable to retreive tenants")
	}

	service := &AddressService{
		repo:                  repo,
		tenantIdToAddressPool: make(map[string]*AddressPool),
	}

	for _, tenant := range tenants {
		addressPool, _ := NewAddressPool(tenant.Network)
		if err != nil {
			log.Fatalf("Unable to instantiate network pool for tenant [%s]", tenant.Id)
		}
		// todo: query for all tenant peers and modify addressPool.used to reflect that the addresses are already in use.
		service.tenantIdToAddressPool[tenant.Id] = addressPool
		log.Printf("Allocating new address pool [%s] for tenant [%s]", tenant.Network, tenant.Id)
	}

	return service
}

func (a *AddressService) AllocateAddress(tenantId string, peerId string, publicKey string) error {

	// Check tenant exist
	tenant, err := a.repo.ReadTenant(tenantId)
	if err != nil {
		return errors.New("Unable to find tenant " + tenantId)
	}

	// Check peer exists on tenant

	// Create new peer object and write to db

	// Add peer entry to health service

	// get all tenant peers

	// send notification peer joined

	// return address response with allocated address to peer

	log.Printf("Allocating network address for peer %s\n", tenant)
	return nil
}

func (a *AddressService) RevokeAddress(tenantId string, peerId string) error {
	tenant, err := a.repo.ReadTenant(tenantId)
	if err != nil {
		return errors.New("Unable to find network for tenant " + tenantId)
	}

	log.Printf("RevokeAddress: tenantId[%s], peer[%s], network[%s]", tenantId, peerId, tenant.Network)
	return nil
}
