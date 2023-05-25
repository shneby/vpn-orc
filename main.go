package main

import (
	"log"
	"vpn-orc/network"
	"vpn-orc/persistence"
)

func main() {
	repo := persistence.NewRepositoryService()
	notifier := network.NewNotificationService()
	addressService := network.NewAddressService(repo, notifier)

	addressService.AllocateAddress(1, "peer1", []byte("asda"), "1.1.1.1")
	addressService.AllocateAddress(1, "peer2", []byte("asdadsadaaa"), "2.2.2.2")
	response, _ := addressService.AllocateAddress(1, "peer3", []byte("asasdsad"), "3.3.3.3")
	log.Println(response)
}
