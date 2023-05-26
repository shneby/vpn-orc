package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"vpn-orc/network"
	"vpn-orc/persistence"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// dependency injection
	repo := persistence.NewRepositoryService()
	notifier := network.NewNotificationService()
	orchestratorService := network.NewOrchestratorService(repo, notifier)
	log.Println(orchestratorService)

	//// simulate work from api
	orchestratorService.AllocateAddress(1, "peer1", []byte("asda"), "1.1.1.1")
	orchestratorService.AllocateAddress(1, "peer2", []byte("asdadsadaaa"), "2.2.2.2")
	orchestratorService.AllocateAddress(1, "peer3", []byte("asasdsad"), "3.3.3.3")

	orchestratorService.AllocateAddress(2, "peer1", []byte("asda"), "4.4.4.4")
	orchestratorService.AllocateAddress(2, "peer2", []byte("asdadsadaaa"), "5.5.5.5")
	orchestratorService.AllocateAddress(2, "peer3", []byte("asasdsad"), "6.6.6.6")

	go func() {
		for range interrupt {
			log.Println("Interrupt received closing...")
			cancel()
		}
	}()
	<-ctx.Done()
}
