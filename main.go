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

	err := addressService.AllocateAddress("tenant1", "peer1", "asdadsad")
	if err != nil {
		log.Printf("Failed to allocate address: %s", err)
	}

	err2 := addressService.RevokeAddress("tenant1", "peer1")
	if err != nil {
		log.Printf("Failed to revoke address: %s", err2)
	}

	//pool, _ := network.NewAddressPool("192.172.16.0/16")
	//
	//for i := 0; i < 1000; i++ {
	//	addr, _ := pool.AllocateAddress()
	//	log.Println(addr)
	//}
	//
	//log.Println("")
}
