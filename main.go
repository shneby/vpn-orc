package main

import (
	"vpn-orc/network"
)

func startGateway() {
	// Just start the controller here it's fine...
}

func main() {

	networkService := network.NewAddressService()
	//notificationService := network.NewNotificationService()
	//startGateway()
	networkService.Main()

}
