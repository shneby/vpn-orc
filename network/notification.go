package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vpn-orc/persistence"
)

type NotificationEvent int64

type NotificationMessage struct {
	EventType NotificationEvent `json:"eventType"`
	Peer      persistence.Peer  `json:"peer"`
}

const (
	Connected NotificationEvent = iota
	Disconnected
)

type NotificationService struct {
}

type NotificationInterface interface {
	NotifyConnected(peer persistence.Peer, peers []persistence.Peer) error
	NotifyDisconnected(peer persistence.Peer, peers []persistence.Peer) error
}

func NewNotificationService() NotificationInterface {
	return &NotificationService{}
}

func (n NotificationService) NotifyConnected(peer persistence.Peer, peers []persistence.Peer) error {
	notification := NotificationMessage{
		EventType: Connected,
		Peer:      peer,
	}
	return notify(notification, peers)
}

func (n NotificationService) NotifyDisconnected(peer persistence.Peer, peers []persistence.Peer) error {
	notification := NotificationMessage{
		EventType: Disconnected,
		Peer:      peer,
	}
	return notify(notification, peers)
}

func notify(notification NotificationMessage, peers []persistence.Peer) error {
	for _, p := range peers {
		url := "https://" + p.RAddr + "/notify"
		marsh, err := json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("failed marshalling notification [%v] - %s", notification, err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(marsh))
		if err != nil {
			return fmt.Errorf("failed building notification [%v] - %s", notification, err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer secret")

		log.Printf("Notify peer %s NotificationMessage sent: %v", p.Id, notification)

		// todo: enable later
		//client := http.Client{Timeout: 3 * time.Second}
		//res, err := client.Do(req)
		//if err != nil {
		//	return fmt.Errorf("failed sending notificaiton request: %s", err)
		//}

		//log.Printf("status Code: %d", res.StatusCode)
	}
	return nil
}
