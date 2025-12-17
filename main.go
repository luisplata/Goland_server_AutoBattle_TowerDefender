package main

import (
	"autobattle-server/game"
	"autobattle-server/network"
	"time"
)

func main() {
	gameManager := game.NewGameManager()
	wsHub := network.NewWsHub()

	httpServer := network.NewHttpServer(gameManager, wsHub)
	go httpServer.Start()

	for {
		games := gameManager.GetAllGames()

		for _, game := range games {
			if game.Clock.ShouldTick() {
				game.Simulation.ProcessTick()

				// 🔥 broadcast snapshot
				snapshot := game.State.GetSnapshot()
				wsHub.Broadcast(game.ID, snapshot)
			}
		}

		time.Sleep(1 * time.Millisecond)
	}
}
