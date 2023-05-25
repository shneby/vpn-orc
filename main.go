package main

import (
	"log"
	"vpn-orc/network"
	"vpn-orc/persistence"
)

func startGateway() {
	// Just start the controller here it's fine...
}

func main() {
	repo := persistence.NewRepositoryService()
	addressService := network.NewAddressService(repo)

	addressService.AllocateAddress("tenant1", "peer1", []byte("asda"))
	addressService.AllocateAddress("tenant1", "peer2", []byte("asdadsadaaa"))
	response, _ := addressService.AllocateAddress("tenant1", "peer3", []byte("asasdsad"))
	log.Println(response)
}
