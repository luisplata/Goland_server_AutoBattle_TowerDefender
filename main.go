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

	lastSnapshots := make(map[int]*game.Snapshot)

	for {
		games := gameManager.GetAllGames()

		for _, g := range games {
			if g.Clock.ShouldTick() {
				g.Simulation.ProcessTick()

				currentSnapshot := game.BuildSnapshot(g.State)

				// Enviar snapshot cada 20 ticks o si es la primera
				if lastSnapshots[g.ID] == nil || g.State.Tick%20 == 0 {
					wsHub.Broadcast(g.ID, currentSnapshot)
					lastSnapshots[g.ID] = &currentSnapshot
				} else {
					// Enviar delta con cambios incrementales
					delta := game.BuildDelta(*lastSnapshots[g.ID], currentSnapshot)
					wsHub.Broadcast(g.ID, delta)
					lastSnapshots[g.ID] = &currentSnapshot
				}
			}
		}

		time.Sleep(1 * time.Millisecond)
	}
}
