package main

import (
	"vpn-orc/network"
	"vpn-orc/persistence"
)

func main() {
	repo := persistence.NewRepositoryService()
	notifier := network.NewNotificationService()
	orchestratorService := network.NewOrchestratorService(repo, notifier)
	gateway := network.NewApiGateway(orchestratorService)
	gateway.ListenAndServe()
}
