package network

import (
	"log"
	"vpn-orc/persistence"
)

type NotificationService struct {
}

type NotificationInterface interface {
	NotifyPeerConnected(peer persistence.Peer, peers []persistence.Peer)
	NotifyPeerDisconnected(peer persistence.Peer, peers []persistence.Peer)
}

func NewNotificationService() NotificationInterface {
	return &NotificationService{}
}

func (n NotificationService) NotifyPeerConnected(peer persistence.Peer, peers []persistence.Peer) {
	log.Printf("Send HTTP request - new peer %s connected, notify %s", peer, peers)
}

func (n NotificationService) NotifyPeerDisconnected(peer persistence.Peer, peers []persistence.Peer) {
	log.Printf("Send HTTP request - disconnected peer %s, notify %s", peer, peers)
}
