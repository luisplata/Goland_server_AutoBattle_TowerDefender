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
				// Guardar la fase anterior antes de procesar
				previousPhase := g.State.GetCurrentPhase()

				g.Simulation.ProcessTick()

				currentSnapshot := game.BuildSnapshot(g.State)

				// Detectar cambio de fase y enviar evento especial
				if g.State.DidPhaseChange() {
					phaseEvent := game.BuildPhaseChangeEvent(g.State, previousPhase)
					wsHub.Broadcast(g.ID, phaseEvent)
				}

				// Detectar cambios en manos y enviar eventos hand_updated
				updatedPlayers := g.State.DrainHandUpdates()
				if len(updatedPlayers) > 0 {
					// Para cada jugador con mano actualizada, enviar evento
					for _, pID := range updatedPlayers {
						if player, ok := currentSnapshot.Players[pID]; ok {
							handEvent := game.BuildHandUpdateEvent(pID, player.Hand, player.DeckCount)
							wsHub.Broadcast(g.ID, handEvent)
						}
					}
				}

				// Enviar snapshot cada tick (máxima actualización)
				wsHub.Broadcast(g.ID, game.SnapshotToUpdate(currentSnapshot))
				lastSnapshots[g.ID] = &currentSnapshot
			}
		}

		time.Sleep(1 * time.Millisecond)
	}
}
