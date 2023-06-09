package network

import (
	"fmt"
	"github.com/robfig/cron"
	_ "github.com/robfig/cron"
	"log"
	"time"
	"vpn-orc/persistence"
)

type OrchestratorService struct {
	repo                  persistence.RepositoryInterface
	notifier              NotificationInterface
	tenantIdToAddressPool map[int]*AddressPool
	tenantIdToPeerIdMap   map[int]map[string]int64
}

type OrchestratorInterface interface {
	OnboardPeer(tenantId int, peerId string, publicKey []byte, peerAddr string) (*persistence.OnboardingResponse, error)
	OffboardPeer(tenantId int, peerId string) error
	UpdatePeer(tenantId int, peerId string)
}

func NewOrchestratorService(repo persistence.RepositoryInterface, notifier NotificationInterface) OrchestratorInterface {
	tenants, err := repo.ReadTenants()
	if err != nil {
		log.Fatalf("unable to retreive tenants from database - %s", err)
	}

	service := &OrchestratorService{
		repo:                  repo,
		notifier:              notifier,
		tenantIdToAddressPool: make(map[int]*AddressPool),
		tenantIdToPeerIdMap:   make(map[int]map[string]int64),
	}

	// create address pools
	// note: in production scenario we would load the network pools state from the database
	for _, tenant := range tenants {
		addressPool, _ := NewAddressPool(tenant.Network)
		if err != nil {
			log.Fatalf("unable to instantiate network pool for tenant [%d]", tenant.Id)
		}
		service.tenantIdToAddressPool[tenant.Id] = addressPool
	}

	// schedule scan for peers to remove every minute
	c := cron.New()
	err = c.AddFunc("0 * * * * *", service.checkPeersToRemove)
	if err != nil {
		log.Fatalf("failed to schedule cron job %s", err)
	}
	c.Start()
	return service
}

func (o *OrchestratorService) OnboardPeer(tenantId int, peerId string, publicKey []byte, peerAddr string) (*persistence.OnboardingResponse, error) {
	log.Printf("received allocation request: tenant [%d], peer [%s]", tenantId, peerId)

	// Check tenant exist
	_, err := o.repo.ReadTenant(tenantId)
	if err != nil {
		return nil, fmt.Errorf("unable to find tenant %d", tenantId)
	}

	// Check peer exists on tenant
	peer, err := o.repo.ReadPeer(tenantId, peerId)
	if peer != nil {
		return nil, fmt.Errorf("peer %s already has an address", peer.Id)
	}

	// Create new peer object and write to db
	networkPool := o.tenantIdToAddressPool[tenantId]
	address, err := networkPool.AllocateAddress()
	if err != nil {
		return nil, fmt.Errorf("unable to allocate addresss for peer %s - %s", peerId, err)
	}

	// get all tenant peers before adding the new peer
	peers, _ := o.repo.ReadPeers(tenantId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch peers for tenant %d - %s", tenantId, err)
	}

	peerDTO := persistence.Peer{Id: peerId, VAddr: address, RAddr: peerAddr, PublicKey: publicKey}
	err = o.repo.WritePeer(tenantId, peerDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to persist peer %s - %s", peer.Id, err)
	}

	// Update in health service & notify
	// Note: Notifications are done synchronously - Given more time we should implement an event system between the orchestrator & notifier
	o.UpdatePeer(tenantId, peerId)
	o.notifier.NotifyConnected(peerDTO, peers)

	// return address response with allocated address to peer
	return &persistence.OnboardingResponse{
		Address: address,
		Peers:   peers,
	}, nil
}

func (o *OrchestratorService) OffboardPeer(tenantId int, peerId string) error {
	// check tenant exists
	tenant, err := o.repo.ReadTenant(tenantId)
	if err != nil {
		return fmt.Errorf("unable to find network for tenant %d", tenantId)
	}

	// fetch peers from the database
	peers, err := o.repo.ReadPeers(tenantId)
	if err != nil {
		return fmt.Errorf("cannot find tenant %d peers", tenantId)
	}

	var peerToRemove persistence.Peer
	peerIndex := -1

	// remove peer from slice
	for i, peer := range peers {
		if peer.Id == peerId {
			peerToRemove = peer
			peerIndex = i
			break
		}
	}

	if peerIndex == -1 {
		return fmt.Errorf("cannot find peer %s on tenant %d", peerId, tenantId)
	}

	peers[peerIndex] = peers[len(peers)-1]
	peers = peers[:len(peers)-1]

	// remove peer from the database
	err = o.repo.DeletePeer(tenantId, peerId)
	if err != nil {
		return fmt.Errorf("cannot delete peer %s for tenant %d", peerId, tenantId)
	}

	// send notification peer removed
	o.notifier.NotifyDisconnected(peerToRemove, peers)

	log.Printf("OffboardPeer: tenantId[%d], peer[%s], network[%s]", tenantId, peerId, tenant.Network)
	return nil
}

func (o *OrchestratorService) UpdatePeer(tenantId int, peerId string) {
	// if peer doesn't exist we should send out a 401 so it can reconnect
	now := time.Now().UnixMilli()
	if o.tenantIdToPeerIdMap[tenantId] == nil {
		o.tenantIdToPeerIdMap[tenantId] = make(map[string]int64)
		o.tenantIdToPeerIdMap[tenantId][peerId] = now
	} else {
		o.tenantIdToPeerIdMap[tenantId][peerId] = time.Now().UnixMilli()
	}
}

func (o *OrchestratorService) checkPeersToRemove() {
	now := time.Now()
	peersToRemove := make(map[int][]string)

	for tenant, peerToTimestamp := range o.tenantIdToPeerIdMap {
		for peerId, timestamp := range peerToTimestamp {
			delta := now.UnixMilli() - timestamp
			if delta > 60000 {
				peersToRemove[tenant] = append(peersToRemove[tenant], peerId)
			}
		}
	}

	for tenantId, peers := range peersToRemove {
		for _, peerId := range peers {
			delete(o.tenantIdToPeerIdMap[tenantId], peerId)
			o.OffboardPeer(tenantId, peerId)
		}
	}
}
